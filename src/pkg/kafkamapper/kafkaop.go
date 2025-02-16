package kafkamapper

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/otterize/intents-operator/src/operator/api/v2beta1"
	"github.com/vishalkuo/bimap"
)

var (
	kafkaOperationToAclOperation = map[v2beta1.KafkaOperation]sarama.AclOperation{
		v2beta1.KafkaOperationAll:             sarama.AclOperationAll,
		v2beta1.KafkaOperationConsume:         sarama.AclOperationRead,
		v2beta1.KafkaOperationProduce:         sarama.AclOperationWrite,
		v2beta1.KafkaOperationCreate:          sarama.AclOperationCreate,
		v2beta1.KafkaOperationDelete:          sarama.AclOperationDelete,
		v2beta1.KafkaOperationAlter:           sarama.AclOperationAlter,
		v2beta1.KafkaOperationDescribe:        sarama.AclOperationDescribe,
		v2beta1.KafkaOperationClusterAction:   sarama.AclOperationClusterAction,
		v2beta1.KafkaOperationDescribeConfigs: sarama.AclOperationDescribeConfigs,
		v2beta1.KafkaOperationAlterConfigs:    sarama.AclOperationAlterConfigs,
		v2beta1.KafkaOperationIdempotentWrite: sarama.AclOperationIdempotentWrite,
	}
	KafkaOperationToAclOperationBMap = bimap.NewBiMapFromMap(kafkaOperationToAclOperation)
)

func KafkaOpFromText(text string) (v2beta1.KafkaOperation, error) {
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
