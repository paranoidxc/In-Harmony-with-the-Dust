<?php
class DefaultCntl extends Cntl
{
    public function befAct()
    {
        echo 'BEFORE ACT' . '<BR>';
    }

    public function actIndex()
    {
        echo 'FFBA' . "<BR>";
        echo 'GET' . "<BR>";
        print_r($this->get);
        echo "<BR>";

        echo 'PSOT' . "<BR>";
        print_r($this->post);

        echo 'render dir output' . "<BR>";
        $this->render('test', array('p1' => 'hello', 'p2' => 'world'), FALSE);

        echo str_repeat("<BR>",5);
        echo 'render return output' . "<BR>";
        $test_out = $this->render('test', array('p1' => 'fuckring', 'p2' => 'world'), TRUE);
        print_r($test_out);
    }


    public function afAct()
    {
        echo 'AFTER ACT' . "<BR>";
    }
}
