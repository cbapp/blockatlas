package notifier

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/trustwallet/blockatlas/db"
	"github.com/trustwallet/blockatlas/pkg/address"
)

const (
	DefaultPushNotificationsBatchLimit = 50

	Notifier = "Notifier"
)

var MaxPushNotificationsBatchLimit uint = DefaultPushNotificationsBatchLimit

func RunNotifier(database *db.Instance, delivery amqp.Delivery) {
	defer func() {
		if err := delivery.Ack(false); err != nil {
			log.WithFields(log.Fields{"service": Notifier}).Error(err)
		}
	}()

	txs, err := GetTransactionsFromDelivery(delivery, Notifier)
	if err != nil {
		log.Error("failed to get transactions", err)
	}

	allAddresses := make([]string, 0)
	for _, tx := range txs {
		allAddresses = append(allAddresses, tx.GetAddresses()...)
	}

	addresses := ToUniqueAddresses(allAddresses)
	for i := range addresses {
		addresses[i] = strconv.Itoa(int(txs[0].Coin)) + "_" + addresses[i]
	}

	fmt.Printf("😋😋😋😋😋😋😋😋😋😋 %+v", addresses)

	if len(txs) < 1 {
		fmt.Println("🐰🐰🐰🐰🐰🐰🐰🐰🐰🐰🐰🐰🐰🐰🐰🐰🐰🐰🐰🐰🐰🐰🐰🐰")
		return
	}
	subscriptionsDataList, err := database.GetSubscriptionsForNotifications(addresses)
	if err != nil || len(subscriptionsDataList) == 0 {
		log.Error("failed to get subscriptionsDataList 🚨🚨", err)
		return
	}

	notifications := make([]TransactionNotification, 0)
	for _, sub := range subscriptionsDataList {
		ua, _, ok := address.UnprefixedAddress(sub.Address.Address)
		if !ok {
			continue
		}
		notificationsForAddress := buildNotificationsByAddress(ua, txs)
		notifications = append(notifications, notificationsForAddress...)
	}

	batches := getNotificationBatches(notifications, MaxPushNotificationsBatchLimit)

	for _, batch := range batches {
		publishNotificationBatch(batch)
	}
	log.Info("------------------------------------------------------------")
}
