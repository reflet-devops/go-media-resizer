import http from 'k6/http';
import { Counter } from 'k6/metrics';
import { check, sleep } from 'k6';
import { selectImageByDistribution } from '../common/utils.js';
import { CONFIG } from '../config/config.js';
const counterByTag = new Counter('counter_by_tag');

function getExpectedResponseTime(width) {
    switch (width) {
        case '1200':
            return CONFIG.scenarios.resize.large.response_time;
        case '800':
            return CONFIG.scenarios.resize.medium.response_time;
        case '400':
            return CONFIG.scenarios.resize.small.response_time;
        default:
            return CONFIG.scenarios.resize.small.response_time;
    }
}

export function resizeTest() {
    const width = __ENV.WIDTH || 800;
    const imagePath = selectImageByDistribution(0, 80, 20);
    const expectedResponseTime = getExpectedResponseTime(width);
    const url = `${CONFIG.baseUrl}${CONFIG.pathPrefix}/${width}${imagePath}`;

    counterByTag.add(1, {test_type: 'resize', width: width});

    const response = http.get(url, {
        tags: {
            test_type: 'resize',
            width: width
        }
    });

    check(response, {
        'resize: status is 200': (r) => r.status === 200,
        'resize: content-type is image': (r) => r.headers['Content-Type'] && r.headers['Content-Type'].startsWith('image/'),
        [`resize: response time < ${expectedResponseTime}ms for width ${width}px`]: (r) => r.timings.duration < expectedResponseTime,
        [`resize: response time < 800ms for width ${width}px`]: (r) => r.timings.duration < 800,
    });

    sleep(1);
}

export const resizeOptions = {
    resize_1200_test: {
        executor: 'constant-arrival-rate',
        rate: CONFIG.scenarios.resize.large.rate,
        timeUnit: CONFIG.scenarios.resize.large.timeUnit,
        preAllocatedVUs: CONFIG.scenarios.resize.large.preAllocatedVUs,
        maxVUs: CONFIG.scenarios.resize.large.maxVUs,
        duration: CONFIG.scenarios.resize.large.duration,
        exec: 'resizeTest',
        env: {
            WIDTH: '1200',
        },
        tags: {test_type: 'resize', width: '1200'},
    },
    resize_800_test: {
        executor: 'constant-arrival-rate',
        rate: CONFIG.scenarios.resize.medium.rate,
        timeUnit: CONFIG.scenarios.resize.medium.timeUnit,
        preAllocatedVUs: CONFIG.scenarios.resize.medium.preAllocatedVUs,
        maxVUs: CONFIG.scenarios.resize.medium.maxVUs,
        duration: CONFIG.scenarios.resize.medium.duration,
        exec: 'resizeTest',
        env: {
            WIDTH: '800',
        },
        tags: {test_type: 'resize', width: '800'}
    },
    resize_400_test: {
        executor: 'constant-arrival-rate',
        rate: CONFIG.scenarios.resize.small.rate,
        timeUnit: CONFIG.scenarios.resize.small.timeUnit,
        preAllocatedVUs: CONFIG.scenarios.resize.small.preAllocatedVUs,
        maxVUs: CONFIG.scenarios.resize.small.maxVUs,
        duration: CONFIG.scenarios.resize.small.duration,
        exec: 'resizeTest',
        env: {
            WIDTH: '400',
        },
        tags: {test_type: 'resize', width: '400'},
    },
};

export const resizeThresholds = {
    'http_req_failed': ['rate<0.01'],
    'http_req_duration{test_type:resize,width:1200}': [`p(95)<${CONFIG.scenarios.resize.large.response_time}`],
    'http_req_duration{test_type:resize,width:800}': [`p(95)<${CONFIG.scenarios.resize.medium.response_time}`],
    'http_req_duration{test_type:resize,width:400}': [`p(95)<${CONFIG.scenarios.resize.small.response_time}`],
    'counter_by_tag{test_type:resize,width:1200}': [],
    'counter_by_tag{test_type:resize,width:800}': [],
    'counter_by_tag{test_type:resize,width:400}': [],
}

export const options = {
    summaryTrendStats: ['min', 'avg', 'med', 'max', 'p(90)', 'p(95)', 'p(99)'],
    scenarios: {
        ...resizeOptions
    },
    thresholds: resizeThresholds
};

export default resizeTest;
