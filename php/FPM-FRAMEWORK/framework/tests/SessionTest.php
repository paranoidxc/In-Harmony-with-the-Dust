<?php
namespace Paranoid\Framework\Tests;

use Paranoid\Framework\Session\Session;
use PHPUnit\Framework\TestCase;

class SessionTest extends TestCase
{
    public function testSetAndGetFlash(): void
    {
        $session = new Session();
        $session->setFlash('success', 'Great job');
        $session->setFlash('error', 'bad job');
        $this->assertTrue($session->hasFlash('success'));
    }

}