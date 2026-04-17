package patterns

// Builtins is the default set of built-in patterns, in match-priority order.
var Builtins = []*Pattern{
	githubPAT,
	githubToken,
	awsAccessKey,
	awsSecretKey,
	stripeKey,
	jwt,
	dbURLPassword,
	pemPrivateKey,
}
