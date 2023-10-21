package kafkamapper

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha3"
	"github.com/vishalkuo/bimap"
)

var (
	kafkaOperationToAclOperation = map[v1alpha3.KafkaOperation]sarama.AclOperation{
		v1alpha3.KafkaOperationAll:             sarama.AclOperationAll,
		v1alpha3.KafkaOperationConsume:         sarama.AclOperationRead,
		v1alpha3.KafkaOperationProduce:         sarama.AclOperationWrite,
		v1alpha3.KafkaOperationCreate:          sarama.AclOperationCreate,
		v1alpha3.KafkaOperationDelete:          sarama.AclOperationDelete,
		v1alpha3.KafkaOperationAlter:           sarama.AclOperationAlter,
		v1alpha3.KafkaOperationDescribe:        sarama.AclOperationDescribe,
		v1alpha3.KafkaOperationClusterAction:   sarama.AclOperationClusterAction,
		v1alpha3.KafkaOperationDescribeConfigs: sarama.AclOperationDescribeConfigs,
		v1alpha3.KafkaOperationAlterConfigs:    sarama.AclOperationAlterConfigs,
		v1alpha3.KafkaOperationIdempotentWrite: sarama.AclOperationIdempotentWrite,
	}
	KafkaOperationToAclOperationBMap = bimap.NewBiMapFromMap(kafkaOperationToAclOperation)
)

func KafkaOpFromText(text string) (v1alpha3.KafkaOperation, error) {
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
