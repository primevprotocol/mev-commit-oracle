package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strconv"

	bidderregistry "github.com/primevprotocol/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primevprotocol/contracts-abi/clients/BlockTracker"
	oracle "github.com/primevprotocol/contracts-abi/clients/Oracle"
	preconfcommitmentstore "github.com/primevprotocol/contracts-abi/clients/PreConfCommitmentStore"
	providerregistry "github.com/primevprotocol/contracts-abi/clients/ProviderRegistry"
	"github.com/primevprotocol/mev-oracle/pkg/events"
	"golang.org/x/sync/errgroup"
)

type BlockStats struct {
	Number                 uint64 `json:"number"`
	Winner                 string `json:"winner"`
	Window                 int64  `json:"window"`
	TotalOpenedCommitments int    `json:"total_opened_commitments"`
	TotalRewards           int    `json:"total_rewards"`
	TotalSlashes           int    `json:"total_slashes"`
	TotalAmount            string `json:"total_amount"`
}

type ProviderBalances struct {
	Provider string `json:"provider"`
	Stake    string `json:"stake"`
	Rewards  string `json:"rewards"`
}

type BidderAllowance struct {
	Bidder    string `json:"bidder"`
	Allowance string `json:"allowance"`
	Refunds   string `json:"refunds"`
	Settled   string `json:"settled"`
	Withdrawn string `json:"withdrawn"`
}

type DashboardOut struct {
	Block     *BlockStats         `json:"block"`
	Providers []*ProviderBalances `json:"providers"`
	Bidders   []*BidderAllowance  `json:"bidders"`
}

func (s *Service) configureDashboard() error {
	blockEvt := events.NewEventHandler(
		"NewL1Block",
		func(upd *blocktracker.BlocktrackerNewL1Block) error {
			s.statMu.Lock()
			defer s.statMu.Unlock()

			existing, ok := s.blockStats.Get(upd.BlockNumber.Uint64())
			if !ok {
				existing = &BlockStats{
					Number: upd.BlockNumber.Uint64(),
				}
			}

			existing.Winner = upd.Winner.Hex()
			existing.Window = upd.Window.Int64()
			_ = s.blockStats.Add(upd.BlockNumber.Uint64(), existing)
			if upd.BlockNumber.Uint64() > s.lastBlock {
				s.lastBlock = upd.BlockNumber.Uint64()
			}
			return nil
		},
	)

	subs := make([]events.Subscription, 0)

	sub, err := s.evtMgr.Subscribe(blockEvt)
	if err != nil {
		return err
	}
	subs = append(subs, sub)

	openedCommitments := events.NewEventHandler(
		"CommitmentStored",
		func(upd *preconfcommitmentstore.PreconfcommitmentstoreCommitmentStored) error {
			s.statMu.Lock()
			defer s.statMu.Unlock()

			existing, ok := s.blockStats.Get(upd.BlockNumber)
			if !ok {
				existing = &BlockStats{
					Number: upd.BlockNumber,
				}
			}

			existing.TotalOpenedCommitments++
			_ = s.blockStats.Add(upd.BlockNumber, existing)
			return nil
		},
	)

	sub, err = s.evtMgr.Subscribe(openedCommitments)
	if err != nil {
		return err
	}
	subs = append(subs, sub)

	settlements := events.NewEventHandler(
		"CommitmentProcessed",
		func(upd *oracle.OracleCommitmentProcessed) error {
			cmt, err := s.store.Settlement(context.Background(), upd.CommitmentHash[:])
			if err != nil {
				return err
			}

			s.statMu.Lock()
			defer s.statMu.Unlock()

			existing, ok := s.blockStats.Get(uint64(cmt.BlockNum))
			if !ok {
				existing = &BlockStats{
					Number: uint64(cmt.BlockNum),
				}
			}

			if upd.IsSlash {
				existing.TotalSlashes++
			} else {
				existing.TotalRewards++
			}
			currentAmount, ok := big.NewInt(0).SetString(existing.TotalAmount, 10)
			if !ok {
				return errors.New("failed to parse total amount")
			}
			currentAmount = big.NewInt(0).Add(currentAmount, big.NewInt(0).SetUint64(cmt.Amount))
			existing.TotalAmount = currentAmount.String()
			_ = s.blockStats.Add(uint64(cmt.BlockNum), existing)
			return nil
		},
	)

	sub, err = s.evtMgr.Subscribe(settlements)
	if err != nil {
		return err
	}
	subs = append(subs, sub)

	providerStakes := events.NewEventHandler(
		"ProviderRegistered",
		func(upd *providerregistry.ProviderregistryProviderRegistered) error {
			s.statMu.Lock()
			defer s.statMu.Unlock()

			existing, ok := s.providerStakes.Get(upd.Provider.Hex())
			if !ok {
				existing = &ProviderBalances{
					Provider: upd.Provider.Hex(),
				}
			}
			existing.Stake = upd.StakedAmount.String()
			_ = s.providerStakes.Add(upd.Provider.Hex(), existing)
			return nil
		},
	)

	sub, err = s.evtMgr.Subscribe(providerStakes)
	if err != nil {
		return err
	}
	subs = append(subs, sub)

	providerDeposit := events.NewEventHandler(
		"FundsDeposited",
		func(upd *providerregistry.ProviderregistryFundsDeposited) error {
			s.statMu.Lock()
			defer s.statMu.Unlock()

			existing, ok := s.providerStakes.Get(upd.Provider.Hex())
			if !ok {
				return errors.New("provider not found")
			}
			currentStake, ok := big.NewInt(0).SetString(existing.Stake, 10)
			if !ok {
				return errors.New("failed to parse stake")
			}
			currentStake = big.NewInt(0).Add(currentStake, upd.Amount)
			existing.Stake = currentStake.String()
			_ = s.providerStakes.Add(upd.Provider.Hex(), existing)
			return nil
		},
	)

	sub, err = s.evtMgr.Subscribe(providerDeposit)
	if err != nil {
		return err
	}
	subs = append(subs, sub)

	providerSlashing := events.NewEventHandler(
		"FundsSlashed",
		func(upd *providerregistry.ProviderregistryFundsSlashed) error {
			s.statMu.Lock()
			defer s.statMu.Unlock()

			existing, ok := s.providerStakes.Get(upd.Provider.Hex())
			if !ok {
				return errors.New("provider not found")
			}
			currentStake, ok := big.NewInt(0).SetString(existing.Stake, 10)
			if !ok {
				return errors.New("failed to parse stake")
			}
			currentStake = big.NewInt(0).Sub(currentStake, upd.Amount)
			existing.Stake = currentStake.String()
			_ = s.providerStakes.Add(upd.Provider.Hex(), existing)
			return nil
		},
	)

	sub, err = s.evtMgr.Subscribe(providerSlashing)
	if err != nil {
		return err
	}
	subs = append(subs, sub)

	providerRewards := events.NewEventHandler(
		"FundsRewarded",
		func(upd *bidderregistry.BidderregistryFundsRewarded) error {
			s.statMu.Lock()
			defer s.statMu.Unlock()

			existing, ok := s.providerStakes.Get(upd.Provider.Hex())
			if !ok {
				return errors.New("provider not found")
			}
			currentRewards, ok := big.NewInt(0).SetString(existing.Rewards, 10)
			if !ok {
				return errors.New("failed to parse rewards")
			}
			currentRewards = big.NewInt(0).Add(currentRewards, upd.Amount)
			existing.Rewards = currentRewards.String()
			_ = s.providerStakes.Add(upd.Provider.Hex(), existing)

			existingBidders, ok := s.bidderAllowances.Get(upd.Window.Uint64())
			if !ok {
				return errors.New("window not found")
			}
			for _, b := range existingBidders {
				if b.Bidder == upd.Bidder.Hex() {
					currentSettled, ok := big.NewInt(0).SetString(b.Settled, 10)
					if !ok {
						return errors.New("failed to parse settled")
					}
					currentSettled = big.NewInt(0).Add(currentSettled, upd.Amount)
					b.Settled = currentSettled.String()
					break
				}
			}
			_ = s.bidderAllowances.Add(upd.Window.Uint64(), existingBidders)

			return nil
		},
	)

	sub, err = s.evtMgr.Subscribe(providerRewards)
	if err != nil {
		return err
	}
	subs = append(subs, sub)

	bidderRegistered := events.NewEventHandler(
		"BidderRegistered",
		func(upd *bidderregistry.BidderregistryBidderRegistered) error {
			s.statMu.Lock()
			defer s.statMu.Unlock()

			existing, ok := s.bidderAllowances.Get(upd.WindowNumber.Uint64())
			if !ok {
				existing = make([]*BidderAllowance, 0)
			}

			for _, b := range existing {
				if b.Bidder == upd.Bidder.Hex() {
					return errors.New("bidder already registered")
				}
			}

			existing = append(existing, &BidderAllowance{
				Bidder:    upd.Bidder.Hex(),
				Allowance: upd.PrepaidAmount.String(),
			})
			_ = s.bidderAllowances.Add(upd.WindowNumber.Uint64(), existing)
			return nil
		},
	)

	sub, err = s.evtMgr.Subscribe(bidderRegistered)
	if err != nil {
		return err
	}
	subs = append(subs, sub)

	bidderPayments := events.NewEventHandler(
		"FundsRetrieved",
		func(upd *bidderregistry.BidderregistryFundsRetrieved) error {
			s.statMu.Lock()
			defer s.statMu.Unlock()

			existing, ok := s.bidderAllowances.Get(upd.Window.Uint64())
			if !ok {
				return errors.New("window not found")
			}

			for _, b := range existing {
				if b.Bidder == upd.Bidder.Hex() {
					currentReturned, ok := big.NewInt(0).SetString(b.Refunds, 10)
					if !ok {
						return errors.New("failed to parse used")
					}
					currentReturned = big.NewInt(0).Add(currentReturned, upd.Amount)
					b.Refunds = currentReturned.String()
					break
				}
			}
			_ = s.bidderAllowances.Add(upd.Window.Uint64(), existing)
			return nil
		},
	)

	sub, err = s.evtMgr.Subscribe(bidderPayments)
	if err != nil {
		return err
	}
	subs = append(subs, sub)

	bidderWithdrawals := events.NewEventHandler(
		"BidderWithdrawal",
		func(upd *bidderregistry.BidderregistryBidderWithdrawal) error {
			s.statMu.Lock()
			defer s.statMu.Unlock()

			existing, ok := s.bidderAllowances.Get(upd.Window.Uint64())
			if !ok {
				return errors.New("window not found")
			}

			for idx, b := range existing {
				if b.Bidder == upd.Bidder.Hex() {
					existing[idx].Withdrawn = upd.Amount.String()
					break
				}
			}

			_ = s.bidderAllowances.Add(upd.Window.Uint64(), existing)
			return nil
		},
	)

	sub, err = s.evtMgr.Subscribe(bidderWithdrawals)
	if err != nil {
		return err
	}
	subs = append(subs, sub)

	eg := errgroup.Group{}
	for _, sub := range subs {
		sub := sub
		eg.Go(func() error {
			select {
			case <-s.shutdown:
				sub.Unsubscribe()
				return nil
			case err := <-sub.Err():
				return err
			}
		})
	}

	closed := make(chan struct{})
	go func() {
		defer close(closed)
		_ = eg.Wait()
	}()

	s.router.Handle("/dashboard", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-closed:
				http.Error(w, "listener closed", http.StatusServiceUnavailable)
			default:
			}

			limit := 10
			limitStr := r.URL.Query().Get("limit")
			if limitStr != "" {
				l, err := strconv.Atoi(limitStr)
				if err == nil {
					limit = l
				}
			}

			page := 0
			pageStr := r.URL.Query().Get("page")
			if pageStr != "" {
				p, err := strconv.Atoi(pageStr)
				if err == nil {
					page = p
				}
			}

			lastBlock := s.lastBlock
			lastBlockStr := r.URL.Query().Get("last_block")
			if lastBlockStr != "" {
				lb, err := strconv.ParseUint(lastBlockStr, 10, 64)
				if err == nil {
					lastBlock = lb
				}
			}

			s.statMu.RLock()
			defer s.statMu.RUnlock()

			start := lastBlock
			if start > uint64(limit*page) {
				start = s.lastBlock - uint64(limit*page)
			}

			dash := make([]*DashboardOut, 0)

			for i := start; i > 0 && len(dash) <= limit; i-- {
				stats, ok := s.blockStats.Get(i)
				if !ok {
					continue
				}
				bidders, ok := s.bidderAllowances.Get(uint64(stats.Window))
				if !ok {
					bidders = make([]*BidderAllowance, 0)
				}

				providers := s.providerStakes.Values()

				dashEntry := &DashboardOut{
					Block:     stats,
					Providers: providers,
					Bidders:   bidders,
				}
				dash = append(dash, dashEntry)
			}

			if err := json.NewEncoder(w).Encode(dash); err != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
		}),
	)

	return nil
}
