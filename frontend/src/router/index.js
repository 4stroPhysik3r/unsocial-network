import { createRouter, createWebHistory } from 'vue-router';
import { useAuth } from '@/composables/useAuth';
import Home from '@/views/Home'
import Groups from '@/views/Groups'
import Chats from '@/views/Chats'
import Profile from '@/views/Profile'
import Register from '@/views/Register'
import Login from '@/views/Login'
import NotFound from '@/views/NotFound'
import User from '@/views/User'
import PostID from '@/views/PostID';
import CreateGroup from '@/views/CreateGroup';
import GroupID from '@/views/GroupID';
import Notifications from '@/views/Notifications';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: { requiresAuth: true }
  },
  {
    path: '/groups',
    name: 'Groups',
    component: Groups,
    meta: { requiresAuth: true }
  },
  {
    path: '/chats',
    name: 'Chats',
    component: Chats,
    meta: { requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: Profile,
    meta: { requiresAuth: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { requiresAuth: false, redirectIfAuthenticated: true } // Add this line
  },
  {
    path: '/login',
    name: "Login",
    component: Login,
    meta: { requiresAuth: false, redirectIfAuthenticated: true } // Add this line
  },
  {
    path: '/userid/:userID', // Define a dynamic route parameter for the userID
    name: 'userid',
    component: User,
    meta: { requiresAuth: true } // Adjust meta as needed
  },

  {
    path: '/postID/:postID', // Define a dynamic route parameter for the postID
    name: 'PostID',
    component: PostID,
    meta: { requiresAuth: true } // Adjust meta as needed
  },
  {
    path: '/create-group',
    name: 'CreateGroup',
    component: CreateGroup,
    meta: { requiresAuth: true }
  },
  // catch all 404 Errors
  {
    path: '/:catchAll(.*)',
    name: 'NotFound',
    component: NotFound
  },
  {
    path: '/groupID/:groupID',
    name: 'GroupID',
    component: GroupID,
    meta: { requiresAuth: true }
  },
  {
    path: '/notifications',
    name: 'Notifications',
    component: Notifications,
    meta: { requiresAuth: true }
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

router.beforeEach(async (to, from, next) => {
  const { isAuthenticated, checkAuthStatus } = useAuth();

  await checkAuthStatus();

  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);

  const redirectIfAuthenticated = to.matched.some(record => record.meta.redirectIfAuthenticated);

  if (requiresAuth && !isAuthenticated.value) {
    next({ name: 'Login' });
  } else if (redirectIfAuthenticated && isAuthenticated.value) {
    next({ name: 'Home' });
  } else {
    next();
  }
});

export default router;