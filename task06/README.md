# Architecture

Below we will describe the communication and IDL schemes for different architecture variants. We will consider schemes for only creating an order.

## RESTful

## IDL schemes

- [Order service](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/arch/restful/order.yml)
- [Billing service](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/arch/restful/billing.yml)
- [Notification service](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/arch/restful/notification.yml)
- [Account service](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/arch/restful/account.yml)

## Communication scheme
![restful](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/images/restful.jpg)

### Notes

Notification service should use a cache to reduce the count of requests to Account service.

# Event Notifications

## IDL schemes

- [Order service](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/arch/event%20notification/order.yml)
- [Billing service](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/arch/event%20notification/billing.yml)
- [Account service](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/arch/event%20notification/account.yml)
- [Async scheme](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/arch/event%20notification/async.yml)

## Communication scheme

![event notification](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/images/event_notification.jpg)

### Notes

Notification service should use a cache to reduce the count of requests to Account service.

## Event Collaboration

All communications are async through a message broker.

## IDL scheme

[Async scheme](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/arch/event%20collaboration/async.yml)

## Communication scheme

![event collaboration](https://github.com/ds-vologdin/otus-software-architect/blob/main/task06/images/event_collaboration.jpg)

### Notes

Notification service should have a view with contacts of users. For updates of this view, the service should listen to events CreateUser and UpdateUser. Here we don't consider these events.
