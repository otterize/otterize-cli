package kafkamapper

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha2"
	"github.com/vishalkuo/bimap"
)

var (
	kafkaOperationToAclOperation = map[v1alpha2.KafkaOperation]sarama.AclOperation{
		v1alpha2.KafkaOperationAll:             sarama.AclOperationAll,
		v1alpha2.KafkaOperationConsume:         sarama.AclOperationRead,
		v1alpha2.KafkaOperationProduce:         sarama.AclOperationWrite,
		v1alpha2.KafkaOperationCreate:          sarama.AclOperationCreate,
		v1alpha2.KafkaOperationDelete:          sarama.AclOperationDelete,
		v1alpha2.KafkaOperationAlter:           sarama.AclOperationAlter,
		v1alpha2.KafkaOperationDescribe:        sarama.AclOperationDescribe,
		v1alpha2.KafkaOperationClusterAction:   sarama.AclOperationClusterAction,
		v1alpha2.KafkaOperationDescribeConfigs: sarama.AclOperationDescribeConfigs,
		v1alpha2.KafkaOperationAlterConfigs:    sarama.AclOperationAlterConfigs,
		v1alpha2.KafkaOperationIdempotentWrite: sarama.AclOperationIdempotentWrite,
	}
	KafkaOperationToAclOperationBMap = bimap.NewBiMapFromMap(kafkaOperationToAclOperation)
)

func KafkaOpFromText(text string) (v1alpha2.KafkaOperation, error) {
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
