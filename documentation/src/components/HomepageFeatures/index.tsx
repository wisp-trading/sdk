import type {ReactNode} from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  image: string;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Type-Safe by Design',
    image: require('@site/static/img/kronos_arrow.png').default,
    description: (
      <>
        Write strategies with full IDE autocomplete and compile-time guarantees.
        Catch errors before they reach production.
      </>
    ),
  },
  {
    title: 'Focus on Strategy Logic',
    image: require('@site/static/img/kronos_target.png').default,
    description: (
      <>
        Don&apos;t worry about exchange APIs or data management.
        Focus on writing your strategies, Kronos handles the rest.
      </>
    ),
  },
  {
    title: 'Write Once, Run Anywhere',
    image: require('@site/static/img/kronos_result.png').default,
    description: (
      <>
        Same code works in backtesting and live trading.
        No environment-specific logic. No adapter layers.
      </>
    ),
  },
];

function Feature({title, image, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <img src={image} className={styles.featureImg} alt={title} />
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
