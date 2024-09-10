<?php
namespace App\Form\User;

use App\Entity\User;
use App\Repository\UserMapper;

class RegistrationForm
{
    private string $username;
    private string $password;
    private array $errors = [];

    public function __construct(
        private UserMapper $userMapper,
    )
    {
    }

    public function setFields(string $username, string $password)
    {
        $this->username = $username;
        $this->password = $password;
    }

    public function save(): User
    {
        $user = User::create($this->username, $this->password);

        $this->userMapper->save($user);

        return $user;
    }

    public function hasValidateionErrors(): bool
    {
       return count($this->getValidationErrors());
    }


    public function getValidationErrors(): array
    {
        if (!empty($this->errors)) {
            return $this->errors;
        }

        if (strlen($this->username)  < 5 || strlen($this->username) > 20) {
            $this->errors[] = 'Username len 5-20';
        }

        if (!preg_match('/^\w+$/', $this->username)) {
            $this->errors[] = "Username work";
        }

        if (strlen($this->password) < 6) {
            $this->errors[] = 'Username len < 6';
        }

        return $this->errors;
    }
}
