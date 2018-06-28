package config

var (
	DefaultSQLDatabaseConfiguration = SQLDatabaseConfiguration{}
	DefaultInfluxDatabaseConfiguration = InfluxDatabaseConfiguration{
		Address: "http://localhost:8086",
		Database: "data",
	}
	DefaultProxyCheckerConfiguration = ProxyCheckerConfiguration{
		StaticUrl: "http://www.example.com",
	}
	DefaultStatisticsReportingConfiguration = StatisticsReportingConfiguration{
		Every: 5000,
		Precision: "s",
	}
	DefaultHttpClientsDefaultConfiguration = HttpClientDefaultsConfiguration{
		UserAgent: "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0",
		Timeout: 1000,
	}
)
