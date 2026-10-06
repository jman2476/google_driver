package auth

// func TestCheckTokenAge(t *testing.T) {
// 	currentTime := time.Now()

// 	dayDuration, err := time.ParseDuration("24h")
// 	if err != nil {
// 		t.Errorf("Error making time duration: %v", err)
// 	}

// 	cases := []struct {
// 		input         oauth2.Token
// 		expectedError error
// 	}{
// 		{
// 			input:         oauth2.Token{Expiry: currentTime.Add(dayDuration)},
// 			expectedError: errTokenNearExpiry,
// 		}, {
// 			input:         oauth2.Token{Expiry: currentTime.Add(-dayDuration)},
// 			expectedError: errTokenExpired,
// 		}, {
// 			input:         oauth2.Token{Expiry: currentTime},
// 			expectedError: errTokenExpired,
// 		}, {
// 			input:         oauth2.Token{Expiry: currentTime.Add(dayDuration).Add(dayDuration)},
// 			expectedError: nil,
// 		},
// 	}

// 	for _, c := range cases {
// 		resultErr := CheckTokenAge(&c.input)
// 		if resultErr != c.expectedError {
// 			t.Errorf("Fail: resulting error does not match expected\nExpected: %v\nActual: %v\n", c.expectedError, resultErr)
// 		}
// 	}
// }
