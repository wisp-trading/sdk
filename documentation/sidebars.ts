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
    {
      type: 'category',
      label: 'Getting Started',
      link: {
        type: 'doc',
        id: 'getting-started/index',
      },
      items: [
        'getting-started/installation',
        'getting-started/quick-reference',
        'getting-started/writing-strategies',
        'getting-started/configuration',
      ],
    },
    {
      type: 'category',
      label: 'Examples',
      link: {
        type: 'doc',
        id: 'examples/index',
      },
      items: [
        {
          type: 'category',
          label: 'Basic Strategies',
          items: [
            'examples/basic/rsi',
            'examples/basic/ma-crossover',
            'examples/basic/bollinger-bands',
          ],
        },
        {
          type: 'category',
          label: 'Intermediate Strategies',
          items: [
            'examples/intermediate/multi-indicator',
            'examples/intermediate/macd-momentum',
            'examples/intermediate/atr-risk',
          ],
        },
        {
          type: 'category',
          label: 'Advanced Strategies',
          items: [
            'examples/advanced/portfolio',
            'examples/advanced/arbitrage',
          ],
        },
      ],
    },
    {
      type: 'category',
      label: 'API Reference',
      items: [
        {
          type: 'category',
          label: 'Indicators',
          items: [
            'api/indicators/moving-averages',
            'api/indicators/rsi',
            'api/indicators/macd',
            'api/indicators/bollinger-bands',
            'api/indicators/stochastic',
            'api/indicators/atr',
          ],
        },
      ],
    },
  ],
};

export default sidebars;