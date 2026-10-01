package ai

const (
	CategoryAIModelsIndustry = "ai_models_industry"
	CategoryAIDevTools       = "ai_dev_tools"
	CategorySecurity         = "security_incident"
	CategorySoftwareDev      = "software_dev"
	CategoryCloudInfra       = "cloud_infra"
	CategoryBusiness         = "business"
	CategoryProductRelease   = "product_release"
	CategoryHardware         = "hardware_gadget"
	CategorySociety          = "society_work"
	CategoryOther            = "other"
)

type ArticleAnalysis struct {
	Category           string
	CategoryConfidence float64
	Importance         float64
}

type ArticleAnalyzer interface {
	Analyze(title, summary, sourceBy string) (ArticleAnalysis, error)
}

var categoryDescriptions = map[string]string{
	CategoryAIModelsIndustry: "AI models and the AI industry: model releases, AI company news, AI safety, regulation, research, and AI-related social issues.",
	CategoryAIDevTools:       "Using AI for development and work: Claude Code, coding agents, MCP, AI-driven development, building apps with AI models such as Jev.",
	CategorySecurity:         "Security and incidents: unauthorized access, data breaches, vulnerabilities, malware, and security products or measures. Includes AI agent security.",
	CategorySoftwareDev:      "Software development where AI is not the main subject: programming languages, frameworks, testing, design, version control, developer tools.",
	CategoryCloudInfra:       "Cloud, infrastructure and networking: AWS, Cloudflare, containers, data centers, telecom networks.",
	CategoryBusiness:         "Corporate and business news: strategy, partnerships, earnings, management, executive interviews, divestitures.",
	CategoryProductRelease:   "Product and service announcements: press releases, new features, new versions, service launches or shutdowns.",
	CategoryHardware:         "Hardware and gadgets: PCs, smartphones, peripherals, chargers, semiconductors, memory.",
	CategorySociety:          "Society, politics, international affairs, incidents, public policy, work styles, remote work, careers.",
	CategoryOther:            "Anything that fits none of the other categories.",
}
