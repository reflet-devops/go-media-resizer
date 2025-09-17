import http from 'k6/http';
import { largeImages, mediumImages, smallImages } from './common/utils.js';
import {CONFIG, SUMMARY_TREND_STATS} from './config/config.js';
import {cdnCgiOptions, cdnCgiTest, cdnCgiThresholds} from './tests/cdn_cgi.js';

// Test scenarios
export let options = {
    scenarios: {
        ...cdnCgiOptions
    },
    summaryTrendStats: SUMMARY_TREND_STATS,
    thresholds: {
        ...cdnCgiThresholds
    },
};

// Setup and teardown functions
export function setup() {
    console.log(`Starting load test against: ${CONFIG.baseUrl}`);
    console.log(`Large images: ${largeImages.length}`);
    console.log(`Medium images: ${mediumImages.length}`);
    console.log(`Small images: ${smallImages.length}`);

    // Test connectivity
    const testResponse = http.get(CONFIG.baseUrl + '/health/ping', {timeout: '5s'});
    if (testResponse.status !== 200 && testResponse.status !== 404) {
        console.warn(`Warning: Base URL ${CONFIG.baseUrl} might not be accessible. Status: ${testResponse.status}`);
    }
}

export function teardown(data) {
    console.log('Load test completed.');
}

export { cdnCgiTest };
