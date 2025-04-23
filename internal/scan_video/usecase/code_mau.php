public function getQuest(Request $request, $computer) {
    $appsetting = Setting::first();
    $dataAva = $this->getAvailableUsernames();
    $channel = $this->getChannel($computer, $dataAva, $appsetting);
    
    if (!$channel) {
        return $this->sendError(null, ['error' => 'account not available']);
    }

    $this->assignProxyIfNeeded($channel);
    $this->updateChannelData($channel);

    $biliSpace = BiliSpace::where('username', $channel->username)->first();
    if (!$biliSpace) {
        return $this->sendError(null, ['error' => 'biliSpace not available']);
    }

    $biliVideo = $this->getBiliVideo($channel, $biliSpace, $appsetting, $computer);
    if (!$biliVideo) {
        $biliSpace->count_upload = $biliSpace->count_video;
        $biliSpace->save();
        $channel->video_left = 1;
        $channel->save();
        return $this->sendError(null, ['error' => $channel->username . ' video bili not available ' . $biliSpace->mid]);
    }

    $this->saveVideoUploading($biliVideo, $channel);
    $this->processVideo($biliVideo, $channel);

    return $this->sendResponse($biliVideo, null);
}

private function getAvailableUsernames() {
    $arryAvai = BiliSpace::select('username')->where('count_video', '>', 'count_upload')->get();
    return $arryAvai->pluck('username')->toArray();
}

private function getChannel($computer, $dataAva, $appsetting) {
    $channel = $computer === 'test_code' 
        ? $this->getTestCodeChannel() 
        : $this->getRegularChannel($computer, $dataAva, $appsetting);

    return $channel;
}

private function getTestCodeChannel() {
    return Channel::where('confirm', '!=', 1)
        ->where('subscribe', '!=', -13)
        ->where('username', 'debojanipathak623@gmail.com')
        ->where('video_left', 0)
        ->where('get_today', '<=', 2)
        ->where('is_upload', 1)
        ->whereNotIn('username', VideoUploading::select('username')->get()->toArray())
        ->where('is_banned', 0)
        ->orderBy('last_get', 'asc')
        ->first();
}

private function getRegularChannel($computer, $dataAva, $appsetting) {
    $currentHour = date('G', time());
    $channel = Channel::where('computer', $computer)
        ->where('confirm', '!=', 1)
        ->where('subscribe', '!=', -13)
        ->where('video_left', 0)
        ->where('is_upload', 1)
        ->where('check_from', '<=', $currentHour)
        ->where('check_to', '>=', $currentHour)
        ->where('upload_today', '=', 0)
        ->where('get_today', '<', 2)
        ->where('is_banned', 0)
        ->whereNotIn('username', VideoUploading::select('username')->get()->toArray())
        ->whereIn('username', $dataAva)
        ->orderBy('last_upload', 'asc')
        ->first();

    if (!$channel) {
        $channel = Channel::where('computer', $computer)
            ->where('confirm', '!=', 1)
            ->where('subscribe', '!=', -13)
            ->where('upload_today', '=', 0)
            ->where('video_left', 0)
            ->where('get_today', '<', 2)
            ->where('is_upload', 1)
            ->where('is_banned', 0)
            ->whereNotIn('username', VideoUploading::select('username')->get()->toArray())
            ->whereIn('username', $dataAva)
            ->orderBy('last_upload', 'asc')
            ->first();
    }

    return $channel;
}

private function assignProxyIfNeeded($channel) {
    if ($channel->proxy == '') {
        $proxies = Proxy::where('live', 1)->orderBy('count_use', 'asc')->take(30)->get()->toArray();

        if (count($proxies) == 0) {
            return $this->sendError(null, ['error' => 'proxy not available']);
        }

        $proxy = $proxies[array_rand($proxies, 1)];
        $channel->proxy = $proxy['proxy_ip'];
        Proxy::where('proxy_ip', $proxy['proxy_ip'])->update(['count_use' => $proxy['count_use'] + 1]);
    }
}

private function updateChannelData($channel) {
    $channel->last_get = time();
    $channel->get_today = $channel->get_today + 1;
    $channel->save();
}

private function getBiliVideo($channel, $biliSpace, $appsetting, $computer) {
    $queCountG = ($channel->next_today > $appsetting->max_next_day) ? '' : 'and count_get < 1';
    return $computer === 'test_code' 
        ? $this->getTestCodeVideo($biliSpace) 
        : $this->getRegularVideo($biliSpace, $queCountG, $channel);
}

private function getTestCodeVideo($biliSpace) {
    return BiliVideo::whereRaw('mid=' . $biliSpace->mid . ' and duration < 1800 and duration !=-1 and video_id not in (select video_id from video_uploading) and video_id not in (select video_id from video_uploaded)')
        ->orderBy('date_make', 'desc')
        ->first();
}

private function getRegularVideo($biliSpace, $queCountG, $channel) {
    if ($channel->new_only == 1) {
        return $this->getNewVideo($biliSpace, $queCountG);
    } else {
        return $this->getAnyVideo($biliSpace, $queCountG);
    }
}

private function getNewVideo($biliSpace, $queCountG) {
    $video = BiliVideo::whereRaw('mid=' . $biliSpace->mid . ' and duration < 1800 ' . $queCountG . ' and download_fail < 3 and next = 0 and duration != -1')
        ->whereNotNull('download_url')
        ->orderBy('date_make', 'desc')
        ->first();

    if (!$video) {
        return BiliVideo::whereRaw('mid=' . $biliSpace->mid . ' and duration < 1800 ' . $queCountG . ' and download_fail < 3 and next = 0 and duration != -1')
            ->where(function($query) {
                $query->whereNull('download_url')->orWhere('download_url', '');
            })
            ->orderBy('date_make', 'desc')
            ->first();
    }

    return $video;
}

private function getAnyVideo($biliSpace, $queCountG) {
    return BiliVideo::whereRaw('mid=' . $biliSpace->mid . ' and duration < 1800 ' . $queCountG . ' and download_fail < 3 and next = 0 and duration != -1')
        ->orderBy('date_make', 'desc')
        ->first();
}

private function saveVideoUploading($biliVideo, $channel) {
    $videoUploading = new VideoUploading();
    $videoUploading->video_id = $biliVideo->video_id;
    $videoUploading->username = $channel->username;
    $videoUploading->start_time = time();
    $videoUploading->mid = $biliSpace->mid;
    $videoUploading->computer = $channel->computer;
    $videoUploading->save();
}

private function processVideo($biliVideo, $channel) {
    if ($biliVideo->duration > 1800) {
        $this->sendError(null, ['error' => 'video max 30 mins']);
    }

    $biliVideo->username = $channel->username;
    $biliVideo->proxy_ip = $channel->proxy;
    $biliVideo->ads = $channel->ads;
    $biliVideo->check_copyright = $channel->check_copyright;

    BiliVideo::where('video_id', $biliVideo->video_id)
        ->update(['duration' => $dataADV['duration'], 'count_get' => ($biliVideo->count_get + 1)]);
}
