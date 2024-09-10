<?php
namespace Paranoid\Framework\Http\Middleware;

use Paranoid\Framework\Http\Request;
use Paranoid\Framework\Http\Response;

interface RequestHandlerInterface
{
    public function handle(Request $request): Response;
}
