const config = JSON.parse(open('./config.json'));

config.baseUrl = __ENV.BASE_URL || config.baseUrl;
config.pathPrefix = __ENV.PATH_PREFIX || config.pathPrefix;

export const SUMMARY_TREND_STATS= ['min', 'avg', 'med', 'max', 'p(90)', 'p(95)', 'p(99)', 'count']
export const CONFIG = config;
