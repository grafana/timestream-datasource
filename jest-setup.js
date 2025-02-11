// Jest setup provided by Grafana scaffolding
import './.config/jest-setup';

Object.defineProperty(global, 'IntersectionObserver', {
  value: jest.fn(() => ({
    observe: jest.fn(),
    unobserve: jest.fn(),
    disconnect: jest.fn(),
    takeRecords: jest.fn(),
    root: null,
    rootMargin: '',
    thresholds: [],
  })),
});
