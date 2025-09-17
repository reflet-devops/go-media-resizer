import http from 'k6/http';
import { Counter } from 'k6/metrics';
import { check, sleep } from 'k6';
import {getRandomElement, scenarioTestEnable, selectImageByDistribution} from '../common/utils.js';
import { CONFIG, SUMMARY_TREND_STATS } from '../config/config.js';
const counterByTag = new Counter('counter_by_tag');

function getExpectedResponseTime(width) {
    switch (width) {
        case '1200':
            return CONFIG.scenarios.resizeFormat.large.response_time;
        case '800':
            return CONFIG.scenarios.resizeFormat.medium.response_time;
        case '400':
            return CONFIG.scenarios.resizeFormat.small.response_time;
        default:
            return CONFIG.scenarios.resizeFormat.small.response_time;
    }
}

export function resizeFormatTest() {
    const width = __ENV.WIDTH || 800;
    const format = getRandomElement(['avif', 'webp']);
    const imagePath = selectImageByDistribution(0, 80, 20);
    const expectedResponseTime = getExpectedResponseTime(width);
    const url = `${CONFIG.baseUrl}${CONFIG.pathPrefix}/${width}${imagePath}`;

    counterByTag.add(1, {test_type: 'resizeFormat', width: width, format: format});

    const response = http.get(url, {
        headers: {
            'Accept': `image/${format},image/png,image/jpeg`,
        },
        tags: {
            test_type: 'resizeFormat',
            width: width,
            format: format
        }
    });

    check(response, {
        'resizeFormat: status is 200': (r) => r.status === 200,
        'resizeFormat: content-type is image': (r) => r.headers['Content-Type'] && r.headers['Content-Type'].startsWith('image/'),
        [`resizeFormat: response time < ${expectedResponseTime}ms for width ${width}px & ${format}`]: (r) => r.timings.duration < expectedResponseTime,
        [`resizeFormat: response time < 800ms for width ${width}px & ${format}`]: (r) => r.timings.duration < 800,
    });

    sleep(1);
}

function generateOptions() {
    let options = {};

    if (scenarioTestEnable(CONFIG.scenarios.resizeFormat.large)) {
        options = {
            ...options,
            resizeFormat_1200_test: {
                executor: 'constant-arrival-rate',
                rate: CONFIG.scenarios.resizeFormat.large.rate,
                timeUnit: CONFIG.scenarios.resizeFormat.large.timeUnit,
                preAllocatedVUs: CONFIG.scenarios.resizeFormat.large.preAllocatedVUs,
                maxVUs: CONFIG.scenarios.resizeFormat.large.maxVUs,
                duration: CONFIG.scenarios.resizeFormat.large.duration,
                exec: 'resizeFormatTest',
                env: {
                    WIDTH: '1200',
                },
                tags: {test_type: 'resizeFormat', width: '1200'},
            }
        };
    }

    if (scenarioTestEnable(CONFIG.scenarios.resizeFormat.medium)) {
        options = {
            ...options,
            resizeFormat_800_test: {
                executor: 'constant-arrival-rate',
                rate: CONFIG.scenarios.resizeFormat.medium.rate,
                timeUnit: CONFIG.scenarios.resizeFormat.medium.timeUnit,
                preAllocatedVUs: CONFIG.scenarios.resizeFormat.medium.preAllocatedVUs,
                maxVUs: CONFIG.scenarios.resizeFormat.medium.maxVUs,
                duration: CONFIG.scenarios.resizeFormat.medium.duration,
                exec: 'resizeFormatTest',
                env: {
                    WIDTH: '800',
                },
                tags: {test_type: 'resizeFormat', width: '800'}
            }
        };
    }

    if (scenarioTestEnable(CONFIG.scenarios.resizeFormat.small)) {
        options = {
            ...options,
            resizeFormat_400_test: {
                executor: 'constant-arrival-rate',
                rate: CONFIG.scenarios.resizeFormat.small.rate,
                timeUnit: CONFIG.scenarios.resizeFormat.small.timeUnit,
                preAllocatedVUs: CONFIG.scenarios.resizeFormat.small.preAllocatedVUs,
                maxVUs: CONFIG.scenarios.resizeFormat.small.maxVUs,
                duration: CONFIG.scenarios.resizeFormat.small.duration,
                exec: 'resizeFormatTest',
                env: {
                    WIDTH: '400',
                },
                tags: {test_type: 'resizeFormat', width: '400'},
            }
        };
    }
    return options;
}

export const resizeFormatOptions = generateOptions();

export const resizeFormatThresholds = {
    'http_req_failed': ['rate<0.01'],
    'counter_by_tag{test_type:resizeFormat,format:avif,width:1200}': [],
    'counter_by_tag{test_type:resizeFormat,format:avif,width:800}': [],
    'counter_by_tag{test_type:resizeFormat,format:avif,width:400}': [],
    'counter_by_tag{test_type:resizeFormat,format:webp,width:1200}': [],
    'counter_by_tag{test_type:resizeFormat,format:webp,width:800}': [],
    'counter_by_tag{test_type:resizeFormat,format:webp,width:400}': [],
    'http_req_duration{test_type:resizeFormat,format:avif,width:1200}': [`p(95)<${CONFIG.scenarios.resizeFormat.large.response_time}`],
    'http_req_duration{test_type:resizeFormat,format:avif,width:800}': [`p(95)<${CONFIG.scenarios.resizeFormat.medium.response_time}`],
    'http_req_duration{test_type:resizeFormat,format:avif,width:400}': [`p(95)<${CONFIG.scenarios.resizeFormat.small.response_time}`],
    'http_req_duration{test_type:resizeFormat,format:webp,width:1200}': [`p(95)<${CONFIG.scenarios.resizeFormat.large.response_time}`],
    'http_req_duration{test_type:resizeFormat,format:webp,width:800}': [`p(95)<${CONFIG.scenarios.resizeFormat.medium.response_time}`],
    'http_req_duration{test_type:resizeFormat,format:webp,width:400}': [`p(95)<${CONFIG.scenarios.resizeFormat.small.response_time}`],
}

export const options = {
    summaryTrendStats: SUMMARY_TREND_STATS,
    scenarios: {
        ...resizeFormatOptions
    },
    thresholds: resizeFormatThresholds,
};

export default resizeFormatTest;
