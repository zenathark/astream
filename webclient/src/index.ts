function audioTest() {
    // var AudioContext = <any>window.AudioContext || <any>window.webkitAudioContext;
    let audioCtx = new AudioContext();
    let buffer = audioCtx.createBuffer(2, audioCtx.sampleRate * 3, audioCtx.sampleRate);
    for (let channel = 0; channel < buffer.numberOfChannels; ++channel) {
	let nowBuffering = buffer.getChannelData(channel);
	for (let i = 0; i < buffer.length; ++i) {
	    nowBuffering[i] = Math.random() * 2 - 1;
	}
    }

    var source = audioCtx.createBufferSource();
    source.buffer = buffer;
    source.connect(audioCtx.destination);
    source.start();
}

document.addEventListener('DOMContentLoaded', audioTest);
