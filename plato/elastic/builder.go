package elastic

type BuilderImpl struct {
}

func NewBuilderImpl() *BuilderImpl {
	return &BuilderImpl{}
}

func (b *BuilderImpl) MatchQuery(term, queryWord string) map[string]interface{} {
	return map[string]interface{}{
		"query": map[string]interface{}{
			"match_phrase": map[string]string{
				term: queryWord,
			},
		},
	}
}

func (b *BuilderImpl) MatchAll() map[string]interface{} {
	return map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
}

func (b *BuilderImpl) MultipleMatch(mappedFields []map[string]string) map[string]interface{} {
	var must []map[string]interface{}

	for _, mappedField := range mappedFields {
		for key, value := range mappedField {
			matchItem := map[string]interface{}{
				"match": map[string]interface{}{
					key: value,
				},
			}
			must = append(must, matchItem)
		}
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
	}

	return query
}

func (b *BuilderImpl) MultiMatchWithGram(queryWord string) map[string]interface{} {
	return map[string]interface{}{
		"size": 15,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": queryWord,
				"type":  "bool_prefix",
				"fields": [3]string{
					"greek", "greek._2gram", "greek._3gram",
				},
			},
		},
	}
}

func (b *BuilderImpl) Aggregate(aggregate, field string) map[string]interface{} {
	return map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			aggregate: map[string]interface{}{
				"terms": map[string]interface{}{
					"field": field,
					"size":  500,
				},
			},
		},
	}
}

func (b *BuilderImpl) FilteredAggregate(term, queryWord, aggregate, field string) map[string]interface{} {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_phrase": map[string]interface{}{
				term: queryWord,
			},
		},
		"size": 0,
		"aggs": map[string]interface{}{
			aggregate: map[string]interface{}{
				"terms": map[string]interface{}{
					"field": field,
					"size":  500,
				},
			},
		},
	}

	return query
}

func (b *BuilderImpl) SearchAsYouTypeIndex(searchWord string) map[string]interface{} {
	return map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				searchWord: map[string]interface{}{
					"type": "search_as_you_type",
				},
			},
		},
	}
}

func (b *BuilderImpl) Index() map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 1,
			},
		},
	}
}
