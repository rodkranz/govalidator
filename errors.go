package govalidator

const (
	// Type mismatch errors.
	ERR_CONTENT_TYPE    = "ContentTypeError"
	ERR_DESERIALIZATION = "DeserializationError"
	ERR_INTERGER_TYPE   = "IntegerTypeError"
	ERR_BOOLEAN_TYPE    = "BooleanTypeError"
	ERR_FLOAT_TYPE      = "FloatTypeError"

	// Validation errors.
	ERR_REQUIRED       = "RequiredError"
	ERR_ALPHA_DASH     = "AlphaDashError"
	ERR_ALPHA_DASH_DOT = "AlphaDashDotError"
	ERR_NUMERIC        = "NumericError"
	ERR_SIZE           = "SizeError"
	ERR_MIN_SIZE       = "MinSizeError"
	ERR_MAX_SIZE       = "MaxSizeError"
	ERR_EMAIL          = "EmailError" //todo
	ERR_URL            = "UrlError"   //todo
	ERR_INCLUDE        = "IncludeError"
	ERR_EXCLUDE        = "ExcludeError"
)

type (
	ErrorHandler interface {
		Error(Errors)
	}

	Validator interface {
		Validate(Errors) Errors
	}

	Errors []Error

	Error struct {
		FieldName      string `json:"field_name,omitempty"`
		FieldAlias     string `json:"field_alias,omitempty"`
		Classification string `json:"classification,omitempty"`
		Message        string `json:"message,omitempty"`
	}
)

func (e *Errors) Add(fieldName, fieldAlias string, classification, message string) {
	*e = append(*e, Error{
		FieldName:      fieldName,
		FieldAlias:     fieldAlias,
		Classification: classification,
		Message:        message,
	})
}

func (e *Errors) Len() int {
	return len(*e)
}

func (e *Errors) Has(class string) bool {
	for _, err := range *e {
		if err.Kind() == class {
			return true
		}
	}
	return false
}


func (e Error) Field() string {
	if len(e.FieldAlias) != 0 {
		return e.FieldAlias
	}
	return e.FieldName
}

func (e Error) Kind() string {
	return e.Classification
}

func (e Error) Error() string {
	return e.Message
}
