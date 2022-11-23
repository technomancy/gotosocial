;; imagine you're tired of your friends talking about pokemon all the time,
;; well, here's the plugin for you!

(fn add-cw? [status]
  (and (= status.ContentWarning "")
       (status.Text:find "pokemon")))

(fn plugin [status]
  (when (add-cw? status)
    (set status.ContentWarning "pokemon")))
