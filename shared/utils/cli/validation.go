package cli

import (
    "fmt"
    "math/big"
    "regexp"
    "strconv"
    "strings"

    "github.com/ethereum/go-ethereum/common"
    "github.com/urfave/cli"

    "github.com/rocket-pool/smartnode/shared/services/passwords"
)


//
// General types
//


// Validate command argument count
func ValidateArgCount(c *cli.Context, count int) error {
    if len(c.Args()) != count {
        return fmt.Errorf("Incorrect argument count; usage: %s", c.Command.UsageText)
    }
    return nil
}


// Validate an address
func ValidateAddress(name, value string) (common.Address, error) {
    if !common.IsHexAddress(value) {
        return common.Address{}, fmt.Errorf("Invalid %s '%s'", name, value)
    }
    return common.HexToAddress(value), nil
}


// Validate a wei amount
func ValidateWeiAmount(name, value string) (*big.Int, error) {
    val := new(big.Int)
    if _, ok := val.SetString(value, 10); !ok {
        return nil, fmt.Errorf("Invalid %s '%s'", name, value)
    }
    return val, nil
}


// Validate a fraction
func ValidateFraction(name, value string) (float64, error) {
    val, err := strconv.ParseFloat(value, 64)
    if err != nil || val < 0 || val > 1 {
        return 0, fmt.Errorf("Invalid %s '%s' - must be a number between 0 and 1", name, value)
    }
    return val, nil
}


// Validate a token type
func ValidateTokenType(name, value string) (string, error) {
    val := strings.ToLower(value)
    if !(val == "eth" || val == "neth") {
        return "", fmt.Errorf("Invalid %s '%s' - valid types are 'ETH' and 'nETH'", name, value)
    }
    return val, nil
}


//
// Command specific types
//


// Validate a deposit amount
func ValidateDepositWeiAmount(name, value string) (*big.Int, error) {
    ether := strings.Repeat("0", 18)
    if !(value == "0" || value == "16"+ether || value == "32"+ether) {
        return nil, fmt.Errorf("Invalid %s '%s' - valid values are 0, 16 and 32 ether", name, value)
    }
    return ValidateWeiAmount(name, value)
}


// Validate a burnable token type
func ValidateBurnableTokenType(name, value string) (string, error) {
    val := strings.ToLower(value)
    if !(val == "neth") {
        return "", fmt.Errorf("Invalid %s '%s' - valid types are 'nETH'", name, value)
    }
    return val, nil
}


// Validate a node password
func ValidateNodePassword(name, value string) (string, error) {
    if len(value) < passwords.MinPasswordLength {
        return "", fmt.Errorf("Invalid %s '%s' - must be at least %d characters long", name, value, passwords.MinPasswordLength)
    }
    return value, nil
}


// Validate a timezone location
func ValidateTimezoneLocation(name, value string) (string, error) {
    if !regexp.MustCompile("^\\w{2,}\\/\\w{2,}$").MatchString(value) {
        return "", fmt.Errorf("Invalid %s '%s' - must be in the format 'Country/City'", name, value)
    }
    return value, nil
}

