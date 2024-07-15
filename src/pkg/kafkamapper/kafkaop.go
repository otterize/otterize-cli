package kafkamapper

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/otterize/intents-operator/src/operator/api/v2alpha1"
	"github.com/vishalkuo/bimap"
)

var (
	kafkaOperationToAclOperation = map[v2alpha1.KafkaOperation]sarama.AclOperation{
		v2alpha1.KafkaOperationAll:             sarama.AclOperationAll,
		v2alpha1.KafkaOperationConsume:         sarama.AclOperationRead,
		v2alpha1.KafkaOperationProduce:         sarama.AclOperationWrite,
		v2alpha1.KafkaOperationCreate:          sarama.AclOperationCreate,
		v2alpha1.KafkaOperationDelete:          sarama.AclOperationDelete,
		v2alpha1.KafkaOperationAlter:           sarama.AclOperationAlter,
		v2alpha1.KafkaOperationDescribe:        sarama.AclOperationDescribe,
		v2alpha1.KafkaOperationClusterAction:   sarama.AclOperationClusterAction,
		v2alpha1.KafkaOperationDescribeConfigs: sarama.AclOperationDescribeConfigs,
		v2alpha1.KafkaOperationAlterConfigs:    sarama.AclOperationAlterConfigs,
		v2alpha1.KafkaOperationIdempotentWrite: sarama.AclOperationIdempotentWrite,
	}
	KafkaOperationToAclOperationBMap = bimap.NewBiMapFromMap(kafkaOperationToAclOperation)
)

func KafkaOpFromText(text string) (v2alpha1.KafkaOperation, error) {
	var saramaOp sarama.AclOperation
	if err := saramaOp.UnmarshalText([]byte(text)); err != nil {
		return "", err
	}

	apiOp, ok := KafkaOperationToAclOperationBMap.GetInverse(saramaOp)
	if !ok {
		return "", fmt.Errorf("failed parsing op %s", saramaOp.String())
	}
	return apiOp, nil
}
