function audioTest() {
    // var AudioContext = <any>window.AudioContext || <any>window.webkitAudioContext;
    var audioCtx = new AudioContext();
    var buffer = audioCtx.createBuffer(2, audioCtx.sampleRate * 3, audioCtx.sampleRate);
    for (var channel = 0; channel < buffer.numberOfChannels; ++channel) {
        var nowBuffering = buffer.getChannelData(channel);
        for (var i = 0; i < buffer.length; ++i) {
            nowBuffering[i] = Math.random() * 2 - 1;
        }
    }
    alert("Ready!");
}
document.addEventListener('DOMContentLoaded', audioTest);
//# sourceMappingURL=index.js.map