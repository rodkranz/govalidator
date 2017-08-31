package govalidator

type (
	Rule struct {
		IsMatch func(string) bool
		IsValid func(Errors, string, string, interface{}) (bool, Errors)
	}

	ParamRule struct {
		IsMatch func(string) bool
		IsValid func(Errors, string, string, string, interface{}) (bool, Errors)
	}

	RuleMapper      []*Rule
	ParamRuleMapper []*ParamRule
)

var ruleMapper RuleMapper
var paramRuleMapper ParamRuleMapper

func AddRule(r ...*Rule) {
	ruleMapper = append(ruleMapper, r...)
}

func AddParamRule(r ...*ParamRule) {
	paramRuleMapper = append(paramRuleMapper, r...)
}
