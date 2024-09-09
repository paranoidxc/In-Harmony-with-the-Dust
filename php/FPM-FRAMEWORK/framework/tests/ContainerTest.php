<?php

namespace Paranoid\Framework\Tests;

use Paranoid\Framework\Container\Container;
use Paranoid\Framework\Http\Request;
use PHPUnit\Framework\TestCase;

class ContainerTest extends TestCase
{
    public function testItWorks()
    {
        $this->assertTrue(true);
    }

    /** @test */
    public function a_service_can_be_retrived_from_the_container()
    {
        // setup
        $container = new Container();
        //$a = new Request([], [], [], [], []);


        // do something
        // id string
        //$container->add('dependant-class', DependantClass::class);

        // make assertions
        //$this->assertInstanceOf(DependantClass::class, $container->get('dependant-class'));
    }

}

