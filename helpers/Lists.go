package helpers

import (
	"github.com/evoila/BPM-Client/model"
)

func MergeStringList(l1, l2 []string) []string {

	if len(l2) > len(l1) {
		return MergeStringList(l2, l1)
	}

	for _, s := range l2 {
		l1 = append(l1, s)
	}

	return l1
}

func MergeMetaDataList(l1, l2 []model.MetaData) []model.MetaData {

	if len(l2) > len(l1) {
		return MergeMetaDataList(l2, l1)
	}

	for _, s := range l2 {
		l1 = append(l1, s)
	}

	return l1
}
