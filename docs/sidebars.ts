import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */
const sidebars: SidebarsConfig = {
  docsSidebar: [
    'intro',
    'quick-start',
    {
      type: 'category',
      label: 'configx - 配置加载',
      link: {
        type: 'doc',
        id: 'modules/configx/overview',
      },
      items: [
        'modules/configx/overview',
        'modules/configx/basic-usage',
        'modules/configx/advanced',
        'modules/configx/api',
      ],
    },
    {
      type: 'category',
      label: 'httpx - HTTP 框架适配器',
      link: {
        type: 'doc',
        id: 'modules/httpx/overview',
      },
      items: [
        'modules/httpx/overview',
        'modules/httpx/usage',
        'modules/httpx/middleware',
        'modules/httpx/huma',
      ],
    },
    {
      type: 'category',
      label: 'logx - 日志记录器',
      link: {
        type: 'doc',
        id: 'modules/logx/overview',
      },
      items: [
        'modules/logx/overview',
        'modules/logx/usage',
        'modules/logx/advanced',
      ],
    },
  ],
};

export default sidebars;
