import http from 'k6/http';
import { Counter } from 'k6/metrics';
import { check, sleep } from 'k6';
import { getRandomElement } from '../common/utils.js';
import { mediumImages } from '../common/utils.js';
import {CONFIG, SUMMARY_TREND_STATS} from '../config/config.js';
const counterByTag = new Counter('counter_by_tag');

export function formatTest() {
    const format = __ENV.FORMAT || 'webp';
    const imagePath = getRandomElement(mediumImages);
    const url = `${CONFIG.baseUrl}${CONFIG.pathPrefix}${imagePath}`;

    counterByTag.add(1, {test_type: 'format', format: format});
    const response = http.get(url, {
        headers: {
            'Accept': `image/${format},image/png,image/jpeg`,
        },
        tags: {
            test_type: 'format',
            format: format
        }
    });

    check(response, {
        'format: status is 200': (r) => r.status === 200,
        'format: content-type is image': (r) => r.headers['Content-Type'] && r.headers['Content-Type'].startsWith('image/'),
        [`format: response time for ${format} < 400ms`]: (r) => r.timings.duration < CONFIG.scenarios.format.response_time,
        [`format: response time for ${format} < 800ms`]: (r) => r.timings.duration < 800,
        'format: correct format returned': (r) => {
            const contentType = r.headers['Content-Type'];
            return contentType && (contentType.includes('avif') || contentType.includes('webp'));
        },
    });

    sleep(1);
}

export const formatOptions = {
    format_avif_test: {
        executor: 'constant-arrival-rate',
        rate: CONFIG.scenarios.format.avif.rate,
        timeUnit: CONFIG.scenarios.format.avif.timeUnit,
        preAllocatedVUs: CONFIG.scenarios.format.avif.preAllocatedVUs,
        maxVUs: CONFIG.scenarios.format.avif.maxVUs,
        duration: CONFIG.scenarios.format.avif.duration,
        exec: 'formatTest',
        env: {
            FORMAT: 'avif',
        },
        tags: {test_type: 'format', format: 'avif'},
    },
    format_webp_test: {
        executor: 'constant-arrival-rate',
        rate: CONFIG.scenarios.format.webp.rate,
        timeUnit: CONFIG.scenarios.format.webp.timeUnit,
        preAllocatedVUs: CONFIG.scenarios.format.webp.preAllocatedVUs,
        maxVUs: CONFIG.scenarios.format.webp.maxVUs,
        duration: CONFIG.scenarios.format.webp.duration,
        exec: 'formatTest',
        env: {
            FORMAT: 'webp',
        },
        tags: {test_type: 'format', format: 'webp'},
    },
};

export const formatThresholds = {
    'http_req_failed': ['rate<0.01'],
    'http_req_duration{test_type:format,format:avif}': [`p(95)<${CONFIG.scenarios.format.response_time}`],
    'http_req_duration{test_type:format,format:webp}': [`p(95)<${CONFIG.scenarios.format.response_time}`],
    'counter_by_tag{test_type:format,format:avif}': [],
    'counter_by_tag{test_type:format,format:webp}': [],
}

export const options = {
    summaryTrendStats: SUMMARY_TREND_STATS,
    scenarios: {
        ...formatOptions
    },
    thresholds: formatThresholds
};

export default formatTest;
