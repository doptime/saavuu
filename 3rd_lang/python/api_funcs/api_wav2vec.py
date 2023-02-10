
from transformers import Wav2Vec2ForCTC, Wav2Vec2CTCTokenizer, AutoFeatureExtractor, Wav2Vec2Model
import torch
import torchaudio
import io
import redis
#https://github.com/pytorch/audio/issues/2363
import librosa
# step2: create a function  to convert ogg to audio vector, using huggingface's Wav2Vec2

wav2vecModel = None
def wav2vecModel():
    global wav2vecModel
    if wav2vecModel == None:
        wav2vecModel = Wav2Vec2Model.from_pretrained("facebook/wav2vec2-base-960h").to('cuda')
    return wav2vecModel
def wav2vec2WithMean(waveform):
    wav2vecModel=wav2vecModel()
    # if waveform is not float32, convert it to float32
    if waveform.dtype != torch.float32:
        waveform = waveform.float()
    # if waveform not in cuda, move it to cuda
    if not waveform.is_cuda:
        waveform = waveform.to('cuda')
    features = wav2vecModel(waveform).last_hidden_state
    features = features.mean(axis=1)
    features = features.mean(axis=0)
    # convert to numpy
    features = features.cpu().detach().numpy().tolist()
    return features


def wav2vec2(waveform):
    wav2vecModel=wav2vecModel()
    # if waveform is not float32, convert it to float32
    if waveform.dtype != torch.float32:
        waveform = waveform.float()
    # if waveform not in cuda, move it to cuda
    if not waveform.is_cuda:
        waveform = waveform.to('cuda')
    features = wav2vecModel(waveform).last_hidden_state
    # convert to numpy
    features = features.cpu().detach().numpy().tolist()
    return features


def api_wav2vec(id, i,send_back):
    # check input her
    if "Data" not in i:
        return send_back(id, {"Err":"missing parameter Data", "Vector": None})
    # your logic here
    audio_data = i["Data"] 
    try:
        mean = i["Mean"] if "Mean" in i else False
        format = i["Format"] if "Format" in i else "mp3"

        audio_file = io.BytesIO(audio_data)
        info = torchaudio.info(io.BytesIO(audio_data))
        if info != None :
            # get audio format,using torchaudio
            format = info.encoding.lower()
        #print("format:", format)
        # reset the file pointer to the beginning
        audio_file.seek(0)

        # load mp3 file
        y, sr = torchaudio.load(audio_file, format=format)
        # convert sampling rate to 16000, using method of resampling
        resampler = torchaudio.transforms.Resample(sr, 16000)
        y = resampler(y)
        if mean:
            vector = wav2vec2WithMean(y)
        else:
            vector = wav2vec2(y)
        send_back(id, {"Vector": vector}, use_single_float=True)
    except Exception as e:
        print("error:", e)
        send_back(id, {"Err": str(e), "Vector": None})

