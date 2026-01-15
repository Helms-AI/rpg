# ngpos-golang-reference-project

## Target Languages

- go

## Types

### GetHealth200Response

GetHealth200Response struct for GetHealth200Response

contains:
- Status: Optional Text

### NullableGetHealth200Response

contains:
- value: Optional GetHealth200Response
- isSet: Boolean

## Functions

### GetStatus

GetStatus returns the Status field value if set, zero value otherwise.

**returns:** Text

**logic:**
```
check for null/empty
return the result
```

### GetStatusOk

GetStatusOk returns a tuple with the Status field value if set, nil otherwise and a boolean to check if the value has been set.

**returns:** Optional string, bool

**logic:**
```
check for null/empty
return the result
```

### HasStatus

HasStatus returns a boolean if a field has been set.

**returns:** Boolean

**logic:**
```
return the result
```

### SetStatus

SetStatus gets a reference to the given string and assigns it to the Status field.

**accepts:**
- v: Text

### MarshalJSON

**returns:** List of byte, error

**logic:**
```
return the result
```

### ToMap

**returns:** Map

### Set

**accepts:**
- val: Optional GetHealth200Response

### IsSet

**returns:** Boolean

**logic:**
```
return the result
```

### Unset

### UnmarshalJSON

**accepts:**
- src: List of byte

**returns:** Error

**logic:**
```
return the result
```

