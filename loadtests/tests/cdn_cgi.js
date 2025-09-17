import http from 'k6/http';
import { Counter } from 'k6/metrics';
import { check, sleep } from 'k6';
import {getRandomElement, selectImageByDistribution, scenarioTestEnable} from '../common/utils.js';
import { CONFIG, SUMMARY_TREND_STATS } from '../config/config.js';
const counterByTag = new Counter('counter_by_tag');

function getExpectedResponseTime(width) {
    switch (width) {
        case '1200':
            return CONFIG.scenarios.cdnCgi.large.response_time;
        case '800':
            return CONFIG.scenarios.cdnCgi.medium.response_time;
        case '400':
            return CONFIG.scenarios.cdnCgi.small.response_time;
        default:
            return CONFIG.scenarios.cdnCgi.small.response_time;
    }
}

export function cdnCgiTest() {
    const width = __ENV.WIDTH || 800;
    const format = getRandomElement(['avif', 'webp']);
    const imagePath = selectImageByDistribution(0, 80, 20);
    const expectedResponseTime = getExpectedResponseTime(width);
    const url = `${CONFIG.baseUrl}/cdn-cgi/image/width=${width},format=auto/${CONFIG.baseUrl}${CONFIG.pathPrefix}${imagePath}`;

    counterByTag.add(1, {test_type: 'cdnCgi', width: width, format: format});

    const response = http.get(url, {
        headers: {
            'Accept': `image/${format},image/png,image/jpeg`,
        },
        tags: {
            test_type: 'cdnCgi',
            width: width,
            format: format
        }
    });

    check(response, {
        'cdnCgi: status is 200': (r) => r.status === 200,
        'cdnCgi: content-type is image': (r) => r.headers['Content-Type'] && r.headers['Content-Type'].startsWith('image/'),
        [`cdnCgi: response time < ${expectedResponseTime}ms for width ${width}px & ${format}`]: (r) => r.timings.duration < expectedResponseTime,
        [`cdnCgi: response time < 800ms for width ${width}px & ${format}`]: (r) => r.timings.duration < 800,
    });

    sleep(1);
}

// export const cdnCgiOptions = {
//     executor: 'constant-arrival-rate',
//     rate: CONFIG.scenarios.cdnCgi.small.rate,
//     timeUnit: CONFIG.scenarios.cdnCgi.small.timeUnit,
//     duration: CONFIG.scenarios.cdnCgi.small.duration,
//     preAllocatedVUs: CONFIG.scenarios.cdnCgi.small.preAllocatedVUs,
//     maxVUs: CONFIG.scenarios.cdnCgi.small.maxVUs,
//     exec: 'cdnCgiTest'
// };

function generateOptions() {
    let options = {};

   if  (scenarioTestEnable(CONFIG.scenarios.cdnCgi.large)) {
       options = {
           ...options,
           cdnCgi_1200_test: {
               executor: 'constant-arrival-rate',
               rate: CONFIG.scenarios.cdnCgi.large.rate,
               timeUnit: CONFIG.scenarios.cdnCgi.large.timeUnit,
               preAllocatedVUs: CONFIG.scenarios.cdnCgi.large.preAllocatedVUs,
               maxVUs: CONFIG.scenarios.cdnCgi.large.maxVUs,
               duration: CONFIG.scenarios.cdnCgi.large.duration,
               exec: 'cdnCgiTest',
               env: {
                   WIDTH: '1200',
               },
               tags: {test_type: 'cdnCgi', width: '1200'},
           },
       };
   }

   if  (scenarioTestEnable(CONFIG.scenarios.cdnCgi.medium)) {
       options = {
           ...options,
           cdnCgi_800_test: {
               executor: 'constant-arrival-rate',
               rate: CONFIG.scenarios.cdnCgi.medium.rate,
               timeUnit: CONFIG.scenarios.cdnCgi.medium.timeUnit,
               preAllocatedVUs: CONFIG.scenarios.cdnCgi.medium.preAllocatedVUs,
               maxVUs: CONFIG.scenarios.cdnCgi.medium.maxVUs,
               duration: CONFIG.scenarios.cdnCgi.medium.duration,
               exec: 'cdnCgiTest',
               env: {
                   WIDTH: '800',
               },
               tags: {test_type: 'cdnCgi', width: '800'}
           },
       };
   }

   if  (scenarioTestEnable(CONFIG.scenarios.cdnCgi.small)) {
       options = {
           ...options,
           cdnCgi_400_test: {
               executor: 'constant-arrival-rate',
               rate: CONFIG.scenarios.cdnCgi.small.rate,
               timeUnit: CONFIG.scenarios.cdnCgi.small.timeUnit,
               preAllocatedVUs: CONFIG.scenarios.cdnCgi.small.preAllocatedVUs,
               maxVUs: CONFIG.scenarios.cdnCgi.small.maxVUs,
               duration: CONFIG.scenarios.cdnCgi.small.duration,
               exec: 'cdnCgiTest',
               env: {
                   WIDTH: '400',
               },
               tags: {test_type: 'cdnCgi', width: '400'},
           }
       };
   }

    return options;
}

export const cdnCgiOptions = generateOptions();

export const cdnCgiThresholds = {
    'http_req_failed': ['rate<0.01'], // Less than 1% failures
    'http_req_duration{test_type:cdnCgi,format:avif,width:1200}': [`p(95)<${CONFIG.scenarios.cdnCgi.large.response_time}`],
    'http_req_duration{test_type:cdnCgi,format:avif,width:800}': [`p(95)<${CONFIG.scenarios.cdnCgi.medium.response_time}`],
    'http_req_duration{test_type:cdnCgi,format:avif,width:400}': [`p(95)<${CONFIG.scenarios.cdnCgi.small.response_time}`],
    'http_req_duration{test_type:cdnCgi,format:webp,width:1200}': [`p(95)<${CONFIG.scenarios.cdnCgi.large.response_time}`],
    'http_req_duration{test_type:cdnCgi,format:webp,width:800}': [`p(95)<${CONFIG.scenarios.cdnCgi.medium.response_time}`],
    'http_req_duration{test_type:cdnCgi,format:webp,width:400}': [`p(95)<${CONFIG.scenarios.cdnCgi.small.response_time}`],
    'counter_by_tag{test_type:cdnCgi}': [],
    'counter_by_tag{test_type:cdnCgi,format:avif,width:1200}': [],
    'counter_by_tag{test_type:cdnCgi,format:avif,width:800}': [],
    'counter_by_tag{test_type:cdnCgi,format:avif,width:400}': [],
    'counter_by_tag{test_type:cdnCgi,format:webp,width:1200}': [],
    'counter_by_tag{test_type:cdnCgi,format:webp,width:800}': [],
    'counter_by_tag{test_type:cdnCgi,format:webp,width:400}': [],
}

export const options = {
    summaryTrendStats: SUMMARY_TREND_STATS,
    scenarios: {
        ...cdnCgiOptions
    },
    thresholds: cdnCgiThresholds
};

export default cdnCgiTest;
