package suplog

func WithFn(fields ...Fields) Fields {
	if len(fields) > 0 && fields[0] != nil {
		result := copyFields(fields[0])
		result["fn"] = DefaultLogger.CallerName()

		return result
	}

	return Fields{
		"fn": DefaultLogger.CallerName(),
	}
}

func WithMore(fields Fields, add Fields) Fields {
	fields = copyFields(fields)
	for k, v := range add {
		fields[k] = v
	}

	return fields
}

func copyFields(fields Fields) Fields {
	ff := make(Fields, len(fields))
	for k, v := range fields {
		ff[k] = v
	}

	return ff
}
