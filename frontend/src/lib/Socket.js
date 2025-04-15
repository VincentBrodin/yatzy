const socket = new WebSocket('ws://localhost:3000/ws');

const emitter = new EventTarget();

socket.addEventListener('open', function(_) {
    console.log("It's open");
});


socket.addEventListener('message', function(event) {
    const data = event.data;

    const encoder = new TextEncoder();
    const buffer = encoder.encode(data).buffer;

    const view = new DataView(buffer);
    const callId = view.getUint32(0);

    const payloadBuffer = buffer.slice(4);
    const payloadView = new DataView(payloadBuffer);

    const decoded = new CustomEvent('message', {
        detail: { callId, message: payloadView }
    });
    emitter.dispatchEvent(decoded);
});

function Send(callId, message) {
    if (socket.readyState <= 1) {
        const encoder = new TextEncoder();
        const msgBytes = encoder.encode(message);

        const buffer = new ArrayBuffer(4 + msgBytes.length);
        const view = new DataView(buffer);

        view.setUint32(0, callId);

        const payload = new Uint8Array(buffer);
        payload.set(msgBytes, 4);

        socket.send(payload);
    }
};


export default {
    Send: Send,
    Emitter: emitter
}
