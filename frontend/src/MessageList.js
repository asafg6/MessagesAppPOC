import React, { Component } from 'react';


class MessageList extends Component {

    constructor(props) {
        super(props);
        this.state = {messages: []}
    }

    componentWillMount() {
        const { events, channel } = this.props;
        events.register(channel, (e) => { 
            this.addMessage(JSON.parse(e.data)); 
        })
    }

    componentWillUnmount() {
        const { events, channel } = this.props;
        events.unRegister(channel, (e) => { this.addMessage(e) })

    }

    addMessage(e) {
        let messages = this.state.messages.slice();
        messages.push(e);
        this.setState({ messages });
    }

    render() {
        const { messages } = this.state; 
        return (
            <div className="message-box" style={{backgroundColor: this.props.color}}>
            { messages.map((message) => {
                return (
                    <div className="msg" key={message.id.toString()} >
                        <p>{message.data}</p>
                    </div>
                );
            }) }
            </div>
        );
    }
}

export default MessageList;