import {SharedArray} from 'k6/data';
import { IMAGES_PATH } from '../config/images.js';

export const largeImages = new SharedArray('large images', () => IMAGES_PATH.large);
export const mediumImages = new SharedArray('medium images', () => IMAGES_PATH.medium);
export const smallImages = new SharedArray('small images', () => IMAGES_PATH.small);


export function getRandomElement(array) {
    return array[Math.floor(Math.random() * array.length)];
}

export function selectImageByDistribution(distribution) {

    if (distribution === "large") {
        return getRandomElement(largeImages);
    } else if (distribution === "medium") {
        return getRandomElement(mediumImages);
    } else {
        return getRandomElement(smallImages);
    }
}

export function getRandomWidth() {
    const widths = ['1200', '800', '400'];
    return getRandomElement(widths);
}

export function scenarioTestEnable(scenarioConfig) {
    if ("enable" in scenarioConfig) {
        return Boolean(scenarioConfig.enable);
    }
    return true;
}
