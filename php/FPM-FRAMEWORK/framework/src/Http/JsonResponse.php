<?php


namespace Paranoid\Framework\Http;

class JsonResponse extends Response
{
    public function __construct(private array $data, private int $code = 200, private string $msg = "操作成功")
    {
        $res = [
            "code" => $code,
            "msg" => $msg,
            "data" => $data,
        ];
        $content = json_encode($res, JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE);
        parent::__construct($content, 200, ["Content-Type: application/json"]);
    }

    public function send(): void
    {
        echo $this->getContent();
        exit;
    }
}