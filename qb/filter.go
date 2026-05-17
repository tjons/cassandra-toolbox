package qb

type operator string

const (
	EqualOperator       operator = "="
	LtOperator          operator = "<"
	GtOperator          operator = ">"
	LteOperator         operator = "<="
	GteOperator         operator = ">="
	NeqOperator         operator = "!="
	InOperator          operator = "IN"
	ContainsOperator    operator = "CONTAINS"
	ContainsKeyOperator operator = "CONTAINS KEY"
)

type filterTerm struct {
	column     string
	operator   operator
	value      any
	values     []any
	deepValues [][]any
}

// Equals is used to compare the column to a single value.
// It will generate a filter term with the "=" operator.
func Equals(value any) filterTerm {
	return filterTerm{
		operator: EqualOperator,
		value:    value,
	}
}

// In is used to compare the column to multiple values.
// It will generate a filter term with the "IN" operator.
func In[T any](values ...T) filterTerm {
	anyValues := make([]any, len(values))
	for i, v := range values {
		anyValues[i] = v
	}

	return filterTerm{
		operator: InOperator,
		values:   anyValues,
	}
}

// CollectionIn is used to compare the column to multiple collection values.
// It will generate a filter term with the "IN" operator.
func CollectionIn(values ...[]any) filterTerm {
	return filterTerm{
		operator:   InOperator,
		deepValues: values,
	}
}

// Contains is used to check if a collection column contains a value.
// It will generate a filter term with the "CONTAINS" operator.
func Contains(value any) filterTerm {
	return filterTerm{
		operator: ContainsOperator,
		value:    value,
	}
}

// ContainsKey is used to check if a map column contains a key.
// It will generate a filter term with the "CONTAINS KEY" operator.
func ContainsKey(value any) filterTerm {
	return filterTerm{
		operator: ContainsKeyOperator,
		value:    value,
	}
}

// LessThan is used to compare the column to a single value using the less than operator.
// It will generate a filter term with the "<" operator.
func LessThan(value any) filterTerm {
	return filterTerm{
		operator: LtOperator,
		value:    value,
	}
}

// GreaterThan is used to compare the column to a single value using the greater than operator.
// It will generate a filter term with the ">" operator.
func GreaterThan(value any) filterTerm {
	return filterTerm{
		operator: GtOperator,
		value:    value,
	}
}

// LessThanEqual is used to compare the column to a single value using the less than or equal operator.
// It will generate a filter term with the "<=" operator.
func LessThanEqual(value any) filterTerm {
	return filterTerm{
		operator: LteOperator,
		value:    value,
	}
}

// GreaterThanEqual is used to compare the column to a single value using the greater than or equal operator.
// It will generate a filter term with the ">=" operator.
func GreaterThanEqual(value any) filterTerm {
	return filterTerm{
		operator: GteOperator,
		value:    value,
	}
}
