import type { Component } from 'solid-js';

import logo from './logo.svg';
import styles from './App.module.css';
import {Posts} from './components/Posts'

const App: Component = () => {
  return (
    <div class={styles.App}>
      <Posts />
    </div>
  );
};

export default App;
