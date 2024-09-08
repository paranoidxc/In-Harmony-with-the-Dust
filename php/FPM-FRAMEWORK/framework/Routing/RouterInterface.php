<?php
namespace Paranoid\Framework\Routing;

use Paranoid\Framework\Http\Request;

interface RouterInterface
{
    public function dispatcher(Request $request);
}
