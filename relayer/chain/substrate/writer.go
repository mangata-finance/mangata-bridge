// Copyright 2020 Snowfork
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"context"
	"encoding/hex"

	"github.com/sirupsen/logrus"

	"golang.org/x/sync/errgroup"

	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/mangata-finance/mangata-bridge/relayer/chain"
)

type Writer struct {
	conn     *Connection
	messages <-chan chain.Message
	log      *logrus.Entry
}

func NewWriter(conn *Connection, messages <-chan chain.Message, log *logrus.Entry) (*Writer, error) {
	return &Writer{
		conn:     conn,
		messages: messages,
		log:      log,
	}, nil
}

func (wr *Writer) Start(ctx context.Context, eg *errgroup.Group) error {
	eg.Go(func() error {
		return wr.writeLoop(ctx)
	})
	return nil
}

func (wr *Writer) writeLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-wr.messages:
			err := wr.Write(ctx, &msg)
			if err != nil {
				wr.log.WithFields(logrus.Fields{
					"appid": hex.EncodeToString(msg.AppID[:]),
					"error": err,
				}).Error("Failure submitting message to substrate")
			}
		}
	}
}

// Write submits a transaction to the chain
func (wr *Writer) Write(_ context.Context, msg *chain.Message) error {

	c, err := types.NewCall(&wr.conn.metadata, "Bridge.submit", msg.AppID, msg.Payload)
	if err != nil {
		return err
	}

	ext := types.NewExtrinsic(c)

	era := types.ExtrinsicEra{IsMortalEra: false}

	genesisHash, err := wr.conn.api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return err
	}

	rv, err := wr.conn.api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return err
	}

	var nonce uint64
	err = wr.conn.api.Client.Call(&nonce, "system_accountNextIndex", "5Ggeh38xrCy49HYrh26QETuQmppaapygTLKQWds7upStaZYK")
	if err != nil {
		return err
	}

	wr.log.WithField("nonce", nonce).Debug("Nonce acquired by rpc.")

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                era,
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: 1,
	}

	extI := ext

	err = extI.Sign(*wr.conn.kp, o)
	if err != nil {
		return err
	}

	_, err = wr.conn.api.RPC.Author.SubmitExtrinsic(extI)
	if err != nil {
		return err
	}

	wr.log.WithFields(logrus.Fields{
		"appid": hex.EncodeToString(msg.AppID[:]),
	}).Info("Submitted message to Substrate")

	return nil
}
