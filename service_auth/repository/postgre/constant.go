package postgre

const CreateUser = `
	INSERT INTO public.users
	(username, email, password, activate_key, is_activated, created_at, updated_at, deleted_at)
	VALUES (:username, :email, :password, :activate_key, :is_activated,:created_at, :updated_at, :deleted_at);
`
const GetUserPinByEmail = `SELECT activate_key FROM public.users WHERE email = $1`

const ActivateUser = `
	UPDATE public.users
	SET is_activated=1, update_at=:update_at
	WHERE email=:email;
`

const GetUser = `SELECT username, email, password FROM public.users WHERE email:email`

const CreateLoginHistory = `INSERT INTO public.login_history
	(username, email, device_id,  created_at, updated_at,  login_at)
	VALUES (:username, :email, :device_id,  :created_at, :updated_at,  :login_at);`

const GetLoginHistories = `SELECT * FROM public.users WHERE email:email LIMIT :limit OFFSET :offset ORDER BY login_at DESC`
