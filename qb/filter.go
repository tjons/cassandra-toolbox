package qb

type Operator string

const (
	EqualOperator       Operator = "="
	LtOperator          Operator = "<"
	GtOperator          Operator = ">"
	LteOperator         Operator = "<="
	GteOperator         Operator = ">="
	NeqOperator         Operator = "!="
	InOperator          Operator = "IN"
	ContainsOperator    Operator = "CONTAINS"
	ContainsKeyOperator Operator = "CONTAINS KEY"
)

type filterTerm struct {
	column     string
	operator   Operator
	value      any
	values     []any
	deepValues [][]any
}

func Equals(value any) filterTerm {
	return filterTerm{
		operator: EqualOperator,
		value:    value,
	}
}

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

func CollectionIn(values ...[]any) filterTerm {
	return filterTerm{
		operator:   InOperator,
		deepValues: values,
	}
}

func Contains(value any) filterTerm {
	return filterTerm{
		operator: ContainsOperator,
		value:    value,
	}
}

func ContainsKey(value any) filterTerm {
	return filterTerm{
		operator: ContainsKeyOperator,
		value:    value,
	}
}

func LessThan(value any) filterTerm {
	return filterTerm{
		operator: LtOperator,
		value:    value,
	}
}

func GreaterThan(value any) filterTerm {
	return filterTerm{
		operator: GtOperator,
		value:    value,
	}
}

func LessThanEqual(value any) filterTerm {
	return filterTerm{
		operator: LteOperator,
		value:    value,
	}
}

func GreaterThanEqual(value any) filterTerm {
	return filterTerm{
		operator: GteOperator,
		value:    value,
	}
}
