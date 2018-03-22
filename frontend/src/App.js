import React, { Component } from 'react';
import Events from './events.js';
import MessageList from './MessageList';
import './App.css';

class App extends Component {

    constructor(props) {
        super(props);
        this.events = new Events();
    }

  render() {
    return (
      <div className="App">

        <div>
            <h1>Messages</h1>
        </div>
        <div style={{display: 'flex', flexDirection: 'row'}}>
            <MessageList channel="red" events={this.events} color="#eb598f" />
            <MessageList channel="yellow" events={this.events} color="#fdcd3b" />
            <MessageList channel="blue" events={this.events} color="#5596f5" />
        </div>
      </div>
    );
  }
}

export default App;