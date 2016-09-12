package e2e

import "sourcegraph.com/sourcegraph/go-selenium"

func init() {
	Register(&Test{
		Name:        "register_flow",
		Description: "Registers a brand new user account via the join page.",
		Func:        testRegisterFlow,
	})
}

func testRegisterFlow(t *T) error {
	// Get join page.
	t.Get(t.Endpoint("/"))
	t.Click(selenium.ByLinkText, "Sign up")

	// Validate username input field.
	t.WaitForElement(selenium.ById, "e2etest-login-field")
	username := t.FindElement(selenium.ById, "e2etest-login-field")
	if username.TagName() != "input" {
		t.Fatalf("username TagName should be input, found %s", username.TagName())
	}
	if username.Text() != "" {
		t.Fatalf("username input field should be empty, found %s", username.Text())
	}
	if !username.IsDisplayed() {
		t.Fatalf("username input field should be displayed")
	}
	if !username.IsEnabled() {
		t.Fatalf("username input field should be enabled")
	}

	// Validate password input field.
	password := t.FindElement(selenium.ById, "e2etest-password-field")
	if password.TagName() != "input" {
		t.Fatalf("password TagName should be input, found %s", password.TagName())
	}
	if password.Text() != "" {
		t.Fatalf("password input field should be empty, found %s", password.Text())
	}
	if !password.IsDisplayed() {
		t.Fatalf("password input field should be displayed")
	}
	if !password.IsEnabled() {
		t.Fatalf("password input field should be enabled")
	}
	if password.IsSelected() {
		t.Fatalf("password input field should not be selected")
	}

	// Validate email input field.
	email := t.FindElement(selenium.ById, "e2etest-email-field")
	if email.TagName() != "input" {
		t.Fatalf("email TagName should be input, found %s", email.TagName())
	}
	if email.Text() != "" {
		t.Fatalf("email input field should be empty, found %s", email.Text())
	}
	if !email.IsDisplayed() {
		t.Fatalf("email input field should be displayed")
	}
	if !email.IsEnabled() {
		t.Fatalf("email input field should be enabled")
	}
	if email.IsSelected() {
		t.Fatalf("email input field should not be selected")
	}

	// Enter username and password for test account.
	username.Click()
	username.SendKeys(t.TestLogin)
	password.Click()
	password.SendKeys("e2etest")
	email.Click()
	email.SendKeys(t.TestEmail)

	// Click the submit button.
	t.Click(selenium.ById, "e2etest-register-button")

	// Wait for redirect to Sourcegraph homepage.
	t.WaitForRedirect(t.Endpoint("/?ob=chrome"), "wait for redirect to home after registration")
	return nil
}
