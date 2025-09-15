import http from 'k6/http';
import { Counter } from 'k6/metrics';
import { check, sleep } from 'k6';
import {scenarioTestEnable, selectImageByDistribution} from '../common/utils.js';
import {CONFIG, SUMMARY_TREND_STATS} from '../config/config.js';
const counterByTag = new Counter('counter_by_tag');

export function sourceTest() {
    const distribution = __ENV.DISTRIBUTION || 'medium';
    const imagePath = selectImageByDistribution(distribution);
    const url = `${CONFIG.baseUrl}${CONFIG.pathPrefix}${imagePath}`;

    counterByTag.add(1, {test_type: 'source', distribution: distribution});
    const response = http.get(url, {
        tags: {test_type: 'source', distribution: distribution}
    });

    check(response, {
        'source: status is 200': (r) => r.status === 200,
        'source: content-type is image': (r) => r.headers['Content-Type'] && r.headers['Content-Type'].startsWith('image/'),
        [`source: response time < ${CONFIG.scenarios.source.response_time}ms with ${distribution}`]: (r) => r.timings.duration < CONFIG.scenarios.source.response_time,
        [`source: response time < 800ms with ${distribution}`]: (r) => r.timings.duration < 800,
    });

    sleep(1);
}

function generateOptions() {
    let options = {};

    if (scenarioTestEnable(CONFIG.scenarios.source.large)) {
        options = {
            ...options,
            source_large_test: {
                executor: 'constant-arrival-rate',
                rate: CONFIG.scenarios.source.large.rate,
                timeUnit: CONFIG.scenarios.source.large.timeUnit,
                preAllocatedVUs: CONFIG.scenarios.source.large.preAllocatedVUs,
                maxVUs: CONFIG.scenarios.source.large.maxVUs,
                duration: CONFIG.scenarios.source.large.duration,
                exec: 'sourceTest',
                env: {
                    DISTRIBUTION: 'large',
                },
                tags: {test_type: 'source', distribution: 'large'},
            }
        };
    }
    if (scenarioTestEnable(CONFIG.scenarios.source.medium)) {
        options = {
            ...options,
            source_medium_test: {
                executor: 'constant-arrival-rate',
                rate: CONFIG.scenarios.source.medium.rate,
                timeUnit: CONFIG.scenarios.source.medium.timeUnit,
                preAllocatedVUs: CONFIG.scenarios.source.medium.preAllocatedVUs,
                maxVUs: CONFIG.scenarios.source.medium.maxVUs,
                duration: CONFIG.scenarios.source.medium.duration,
                exec: 'sourceTest',
                env: {
                    DISTRIBUTION: 'medium',
                },
                tags: {test_type: 'source', distribution: 'medium'},
            }
        };
    }

    if (scenarioTestEnable(CONFIG.scenarios.source.small)) {
        options = {
            ...options,
            source_small_test: {
                executor: 'constant-arrival-rate',
                rate: CONFIG.scenarios.source.small.rate,
                timeUnit: CONFIG.scenarios.source.small.timeUnit,
                preAllocatedVUs: CONFIG.scenarios.source.small.preAllocatedVUs,
                maxVUs: CONFIG.scenarios.source.small.maxVUs,
                duration: CONFIG.scenarios.source.small.duration,
                exec: 'sourceTest',
                env: {
                    DISTRIBUTION: 'small',
                },
                tags: {test_type: 'source', distribution: 'small'},
            }
        };
    }
    return options;
}

export const sourceOptions = generateOptions();

export const sourceThresholds = {
    'http_req_failed': ['rate<0.01'],
    'http_req_duration{test_type:source,distribution:small}': [`p(95)<${CONFIG.scenarios.source.response_time}`],
    'http_req_duration{test_type:source,distribution:medium}': [`p(95)<${CONFIG.scenarios.source.response_time}`],
    'http_req_duration{test_type:source,distribution:large}': [`p(95)<${CONFIG.scenarios.source.response_time}`],
    'counter_by_tag{test_type:source,distribution:small}': [],
    'counter_by_tag{test_type:source,distribution:medium}': [],
    'counter_by_tag{test_type:source,distribution:large}': [],
}

export const options = {
    summaryTrendStats: SUMMARY_TREND_STATS,
    scenarios: {
        ...sourceOptions
    },
    thresholds: sourceThresholds
};

export default sourceTest;
