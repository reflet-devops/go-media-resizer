import http from 'k6/http';
import {largeImages, mediumImages, smallImages} from './common/utils.js';
import {CONFIG, SUMMARY_TREND_STATS} from './config/config.js';
import {sourceOptions, sourceTest, sourceThresholds} from './tests/source.js';
import {formatOptions, formatTest, formatThresholds} from './tests/format.js';
import {resizeOptions, resizeTest, resizeThresholds} from './tests/resize.js';
import {resizeFormatOptions, resizeFormatTest, resizeFormatThresholds} from './tests/resize_format.js';

// Test scenarios
export let options = {
    scenarios: {
        ...sourceOptions,
        ...formatOptions,
        ...resizeOptions,
        ...resizeFormatOptions,

    },
    summaryTrendStats: SUMMARY_TREND_STATS,
    thresholds: {
        ...sourceThresholds,
        ...formatThresholds,
        ...resizeThresholds,
        ...resizeFormatThresholds,
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

export {sourceTest, formatTest, resizeTest, resizeFormatTest};
