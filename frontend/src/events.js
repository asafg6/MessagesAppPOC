

class Events {

    constructor() {
        this.client = new EventSource('https://' + window.location.host + '/messages');
    }

    register(eventName, handler) {
        this.client.addEventListener(eventName, handler);
    }

    unRegister(eventName, handler) {
        this.client.removeEventListener(eventName, handler);
    }


}

export default Events;