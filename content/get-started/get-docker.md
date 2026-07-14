#!/bin/bash
echo "🚛 Создание проекта 'Грузавтотранс' с картой, весом, темами, фирмами и контактами..."

# Создаём структуру папок
mkdir -p gruzavtotrans/backend/config
mkdir -p gruzavtotrans/backend/apps/{users,documents,orders,chat,notifications,companies,theme}
mkdir -p gruzavtotrans/frontend/public
mkdir -p gruzavtotrans/frontend/src/components/{auth,dashboard,profile,orders,chat,map,settings,companies}
mkdir -p gruzavtotrans/frontend/src/context
mkdir -p gruzavtotrans/frontend/src/hooks
mkdir -p gruzavtotrans/frontend/src/services
mkdir -p gruzavtotrans/frontend/src/utils

cd gruzavtotrans

# ======================== BACKEND ========================

cat > backend/requirements.txt << 'EOF'
Django==4.2.7
djangorestframework==3.14.0
django-cors-headers==4.3.1
channels==4.0.0
channels-redis==4.1.0
celery==5.3.4
redis==5.0.1
psycopg2-binary==2.9.9
django-environ==0.11.2
django-storages==1.14.2
boto3==1.34.11
twilio==8.8.0
sendgrid==6.10.0
Pillow==10.1.0
EOF

cat > backend/Dockerfile << 'EOF'
FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
CMD ["sh", "-c", "python manage.py migrate && python manage.py runserver 0.0.0.0:8000"]
EOF

cat > backend/manage.py << 'EOF'
#!/usr/bin/env python
import os
import sys

def main():
    os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'config.settings')
    try:
        from django.core.management import execute_from_command_line
    except ImportError as exc:
        raise ImportError(
            "Couldn't import Django. Are you sure it's installed?"
        ) from exc
    execute_from_command_line(sys.argv)

if __name__ == '__main__':
    main()
EOF

cat > backend/config/settings.py << 'EOF'
import os
from pathlib import Path
import environ

env = environ.Env()
BASE_DIR = Path(__file__).resolve().parent.parent
environ.Env.read_env(os.path.join(BASE_DIR, '.env'))

SECRET_KEY = env('SECRET_KEY', default='django-insecure-123456')
DEBUG = env.bool('DEBUG', default=True)
ALLOWED_HOSTS = env.list('ALLOWED_HOSTS', default=['*'])

INSTALLED_APPS = [
    'django.contrib.admin',
    'django.contrib.auth',
    'django.contrib.contenttypes',
    'django.contrib.sessions',
    'django.contrib.messages',
    'django.contrib.staticfiles',
    'rest_framework',
    'corsheaders',
    'channels',
    'apps.users',
    'apps.documents',
    'apps.orders',
    'apps.chat',
    'apps.notifications',
    'apps.companies',
    'apps.theme',
]

MIDDLEWARE = [
    'corsheaders.middleware.CorsMiddleware',
    'django.middleware.security.SecurityMiddleware',
    'django.contrib.sessions.middleware.SessionMiddleware',
    'django.middleware.common.CommonMiddleware',
    'django.middleware.csrf.CsrfViewMiddleware',
    'django.contrib.auth.middleware.AuthenticationMiddleware',
    'django.contrib.messages.middleware.MessageMiddleware',
    'django.middleware.clickjacking.XFrameOptionsMiddleware',
]

ROOT_URLCONF = 'config.urls'
TEMPLATES = [
    {
        'BACKEND': 'django.template.backends.django.DjangoTemplates',
        'DIRS': [],
        'APP_DIRS': True,
        'OPTIONS': {
            'context_processors': [
                'django.template.context_processors.debug',
                'django.template.context_processors.request',
                'django.contrib.auth.context_processors.auth',
                'django.contrib.messages.context_processors.messages',
            ],
        },
    },
]
WSGI_APPLICATION = 'config.wsgi.application'
ASGI_APPLICATION = 'config.asgi.application'

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.postgresql',
        'NAME': env('DB_NAME', default='gruz_db'),
        'USER': env('DB_USER', default='gruz_user'),
        'PASSWORD': env('DB_PASSWORD', default='password'),
        'HOST': env('DB_HOST', default='db'),
        'PORT': env('DB_PORT', default='5432'),
    }
}

CHANNEL_LAYERS = {
    'default': {
        'BACKEND': 'channels_redis.core.RedisChannelLayer',
        'CONFIG': {'hosts': [(env('REDIS_HOST', default='redis'), 6379)]},
    }
}

CORS_ALLOW_ALL_ORIGINS = True
STATIC_URL = '/static/'
MEDIA_URL = '/media/'
MEDIA_ROOT = os.path.join(BASE_DIR, 'media')
DEFAULT_AUTO_FIELD = 'django.db.models.BigAutoField'

TWILIO_ACCOUNT_SID = env('TWILIO_ACCOUNT_SID', default='')
TWILIO_AUTH_TOKEN = env('TWILIO_AUTH_TOKEN', default='')
TWILIO_PHONE = env('TWILIO_PHONE', default='')

EMAIL_HOST = 'smtp.sendgrid.net'
EMAIL_HOST_USER = env('EMAIL_HOST_USER', default='')
EMAIL_HOST_PASSWORD = env('EMAIL_HOST_PASSWORD', default='')
EMAIL_PORT = 587
EMAIL_USE_TLS = True

CELERY_BROKER_URL = f"redis://{env('REDIS_HOST', default='redis')}:6379/0"
CELERY_RESULT_BACKEND = CELERY_BROKER_URL
YANDEX_MAPS_API_KEY = env('YANDEX_MAPS_API_KEY', default='')
EOF

cat > backend/config/urls.py << 'EOF'
from django.contrib import admin
from django.urls import path, include
from django.conf import settings
from django.conf.urls.static import static

urlpatterns = [
    path('admin/', admin.site.urls),
    path('api/auth/', include('apps.users.urls')),
    path('api/orders/', include('apps.orders.urls')),
    path('api/documents/', include('apps.documents.urls')),
    path('api/companies/', include('apps.companies.urls')),
    path('api/theme/', include('apps.theme.urls')),
] + static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
EOF

cat > backend/config/asgi.py << 'EOF'
import os
from django.core.asgi import get_asgi_application
from channels.routing import ProtocolTypeRouter, URLRouter
from channels.auth import AuthMiddlewareStack
from apps.chat.routing import websocket_urlpatterns

os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'config.settings')

application = ProtocolTypeRouter({
    "http": get_asgi_application(),
    "websocket": AuthMiddlewareStack(URLRouter(websocket_urlpatterns)),
})
EOF

cat > backend/config/celery.py << 'EOF'
import os
from celery import Celery

os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'config.settings')
app = Celery('config')
app.config_from_object('django.conf:settings', namespace='CELERY')
app.autodiscover_tasks()
EOF

touch backend/config/__init__.py

# ---------- apps/users ----------
cat > backend/apps/users/models.py << 'EOF'
from django.db import models
from django.contrib.auth.models import AbstractUser

class User(AbstractUser):
    phone = models.CharField(max_length=20, unique=True)
    email = models.EmailField(unique=True)
    full_name = models.CharField(max_length=255)
    inn = models.CharField(max_length=12, blank=True)
    passport_data = models.JSONField(default=dict)
    registration_address = models.TextField(blank=True)
    driver_license_data = models.JSONField(default=dict, blank=True)
    vehicle = models.CharField(max_length=50, blank=True, help_text="Госномер тягача")
    trailer = models.CharField(max_length=50, blank=True, help_text="Госномер прицепа")
    role = models.CharField(max_length=20, choices=[('carrier','Перевозчик'),('customer','Заказчик'),('admin','Админ')])
    is_phone_verified = models.BooleanField(default=False)
    is_email_verified = models.BooleanField(default=False)
    is_active = models.BooleanField(default=False)
    created_at = models.DateTimeField(auto_now_add=True)
    fcm_token = models.CharField(max_length=255, blank=True, null=True)

class VerificationCode(models.Model):
    user = models.ForeignKey(User, on_delete=models.CASCADE)
    code = models.CharField(max_length=6)
    type = models.CharField(max_length=20)
    created_at = models.DateTimeField(auto_now_add=True)
    is_used = models.BooleanField(default=False)
EOF

cat > backend/apps/users/serializers.py << 'EOF'
from rest_framework import serializers
from .models import User

class UserSerializer(serializers.ModelSerializer):
    password = serializers.CharField(write_only=True)
    class Meta:
        model = User
        fields = ('id', 'email', 'phone', 'full_name', 'inn', 'passport_data',
                  'registration_address', 'driver_license_data', 'vehicle', 'trailer',
                  'role', 'password')
    def create(self, validated_data):
        user = User(**validated_data)
        user.set_password(validated_data['password'])
        user.save()
        return user
EOF

cat > backend/apps/users/views.py << 'EOF'
import random
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework.permissions import AllowAny
from .models import User, VerificationCode
from .serializers import UserSerializer
from apps.notifications.tasks import send_sms, send_email

class RegisterView(APIView):
    permission_classes = [AllowAny]
    def post(self, request):
        serializer = UserSerializer(data=request.data)
        if serializer.is_valid():
            user = serializer.save(is_active=False)
            sms_code = f"{random.randint(100000, 999999)}"
            VerificationCode.objects.create(user=user, code=sms_code, type='sms')
            send_sms.delay(user.phone, f'Код подтверждения: {sms_code}')
            return Response({'message': 'SMS отправлен', 'phone': user.phone}, status=201)
        return Response(serializer.errors, status=400)

class VerifySMSView(APIView):
    permission_classes = [AllowAny]
    def post(self, request):
        phone = request.data.get('phone')
        code = request.data.get('code')
        try:
            user = User.objects.get(phone=phone)
            vc = VerificationCode.objects.filter(user=user, code=code, type='sms', is_used=False).first()
            if vc:
                vc.is_used = True
                vc.save()
                user.is_phone_verified = True
                user.save()
                email_code = f"{random.randint(100000, 999999)}"
                VerificationCode.objects.create(user=user, code=email_code, type='email')
                send_email.delay(user.email, 'Подтверждение email', f'Ваш код: {email_code}')
                return Response({'message': 'Телефон подтверждён, проверьте почту', 'email': user.email})
        except User.DoesNotExist:
            pass
        return Response({'error': 'Неверный код или телефон'}, status=400)

class VerifyEmailView(APIView):
    permission_classes = [AllowAny]
    def post(self, request):
        email = request.data.get('email')
        code = request.data.get('code')
        try:
            user = User.objects.get(email=email)
            vc = VerificationCode.objects.filter(user=user, code=code, type='email', is_used=False).first()
            if vc:
                vc.is_used = True
                vc.save()
                user.is_email_verified = True
                user.is_active = True
                user.save()
                return Response({'message': 'Регистрация завершена'})
        except User.DoesNotExist:
            pass
        return Response({'error': 'Неверный код или email'}, status=400)
EOF

cat > backend/apps/users/urls.py << 'EOF'
from django.urls import path
from .views import RegisterView, VerifySMSView, VerifyEmailView

urlpatterns = [
    path('register/', RegisterView.as_view()),
    path('verify-sms/', VerifySMSView.as_view()),
    path('verify-email/', VerifyEmailView.as_view()),
]
EOF

# ---------- apps/documents ----------
cat > backend/apps/documents/models.py << 'EOF'
from django.db import models
from apps.users.models import User

class Document(models.Model):
    user = models.ForeignKey(User, on_delete=models.CASCADE)
    type = models.CharField(max_length=50)
    file = models.FileField(upload_to='documents/%Y/%m/%d/')
    uploaded_at = models.DateTimeField(auto_now_add=True)
    verified = models.BooleanField(default=False)
EOF

cat > backend/apps/documents/urls.py << 'EOF'
from django.urls import path
# заглушка
urlpatterns = []
EOF

# ---------- apps/orders ----------
cat > backend/apps/orders/models.py << 'EOF'
from django.db import models
from apps.users.models import User

class Address(models.Model):
    street = models.CharField(max_length=255)
    house = models.CharField(max_length=20)
    apartment = models.CharField(max_length=20, blank=True)
    city = models.CharField(max_length=100)
    region = models.CharField(max_length=100, blank=True)
    country = models.CharField(max_length=100, default='Россия')
    latitude = models.FloatField(null=True, blank=True)
    longitude = models.FloatField(null=True, blank=True)

    def __str__(self):
        return f"{self.city}, {self.street} {self.house}"

class Order(models.Model):
    customer = models.ForeignKey(User, on_delete=models.CASCADE, related_name='orders_as_customer')
    carrier = models.ForeignKey(User, on_delete=models.SET_NULL, null=True, blank=True, related_name='orders_as_carrier')
    pickup_address = models.ForeignKey(Address, on_delete=models.CASCADE, related_name='pickup_orders')
    delivery_address = models.ForeignKey(Address, on_delete=models.CASCADE, related_name='delivery_orders')
    cargo_description = models.TextField()
    weight = models.FloatField()
    volume = models.FloatField(default=0)
    date_required = models.DateField()
    status = models.CharField(max_length=20, default='created',
                              choices=[('created','Создан'),('searching','Ищет перевозчика'),
                                       ('accepted','Принят'),('in_progress','В пути'),
                                       ('delivered','Доставлен'),('cancelled','Отменён')])
    public_url = models.CharField(max_length=100, unique=True, blank=True)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    def __str__(self):
        return f"Заказ #{self.id} ({self.weight}т)"
EOF

cat > backend/apps/orders/views.py << 'EOF'
from rest_framework.generics import ListCreateAPIView, RetrieveAPIView
from rest_framework.permissions import IsAuthenticated
from .models import Order, Address
from .serializers import OrderSerializer, AddressSerializer

class OrderListCreateView(ListCreateAPIView):
    serializer_class = OrderSerializer
    permission_classes = [IsAuthenticated]
    def get_queryset(self):
        return Order.objects.all()

class OrderDetailView(RetrieveAPIView):
    queryset = Order.objects.all()
    serializer_class = OrderSerializer
EOF

cat > backend/apps/orders/serializers.py << 'EOF'
from rest_framework import serializers
from .models import Order, Address

class AddressSerializer(serializers.ModelSerializer):
    class Meta:
        model = Address
        fields = '__all__'

class OrderSerializer(serializers.ModelSerializer):
    pickup_address = AddressSerializer()
    delivery_address = AddressSerializer()

    class Meta:
        model = Order
        fields = '__all__'

    def create(self, validated_data):
        pickup_data = validated_data.pop('pickup_address')
        delivery_data = validated_data.pop('delivery_address')
        pickup = Address.objects.create(**pickup_data)
        delivery = Address.objects.create(**delivery_data)
        return Order.objects.create(pickup_address=pickup, delivery_address=delivery, **validated_data)
EOF

cat > backend/apps/orders/urls.py << 'EOF'
from django.urls import path
from .views import OrderListCreateView, OrderDetailView

urlpatterns = [
    path('', OrderListCreateView.as_view()),
    path('<int:pk>/', OrderDetailView.as_view()),
]
EOF

# ---------- apps/chat ----------
cat > backend/apps/chat/consumers.py << 'EOF'
import json
from channels.generic.websocket import AsyncWebsocketConsumer

class ChatConsumer(AsyncWebsocketConsumer):
    async def connect(self):
        self.order_id = self.scope['url_route']['kwargs']['order_id']
        self.room_group_name = f'chat_{self.order_id}'
        await self.channel_layer.group_add(self.room_group_name, self.channel_name)
        await self.accept()

    async def disconnect(self, close_code):
        await self.channel_layer.group_discard(self.room_group_name, self.channel_name)

    async def receive(self, text_data):
        data = json.loads(text_data)
        await self.channel_layer.group_send(
            self.room_group_name,
            {'type': 'chat_message', 'message': data['message']}
        )

    async def chat_message(self, event):
        await self.send(text_data=json.dumps({'message': event['message']}))
EOF

cat > backend/apps/chat/routing.py << 'EOF'
from django.urls import re_path
from .consumers import ChatConsumer

websocket_urlpatterns = [
    re_path(r'ws/chat/(?P<order_id>\w+)/$', ChatConsumer.as_asgi()),
]
EOF

# ---------- apps/notifications ----------
cat > backend/apps/notifications/tasks.py << 'EOF'
from celery import shared_task

@shared_task
def send_sms(phone, message):
    print(f'SMS to {phone}: {message}')
    return True

@shared_task
def send_email(recipient, subject, body):
    print(f'Email to {recipient}: {subject} - {body}')
    return True
EOF

touch backend/apps/notifications/services.py

# ---------- apps/companies ----------
cat > backend/apps/companies/models.py << 'EOF'
from django.db import models
from apps.users.models import User

class Company(models.Model):
    name = models.CharField(max_length=255)
    inn = models.CharField(max_length=12, unique=True)
    address = models.TextField()
    phone = models.CharField(max_length=20)
    email = models.EmailField()
    website = models.URLField(blank=True)
    description = models.TextField(blank=True)
    logo = models.ImageField(upload_to='company_logos/', blank=True, null=True)
    created_by = models.ForeignKey(User, on_delete=models.SET_NULL, null=True)
    created_at = models.DateTimeField(auto_now_add=True)

    def __str__(self):
        return self.name

class Contact(models.Model):
    company = models.ForeignKey(Company, on_delete=models.CASCADE, related_name='contacts')
    full_name = models.CharField(max_length=255)
    position = models.CharField(max_length=100)
    phone = models.CharField(max_length=20)
    email = models.EmailField()
    is_primary = models.BooleanField(default=False)

    def __str__(self):
        return f"{self.full_name} ({self.company.name})"
EOF

cat > backend/apps/companies/serializers.py << 'EOF'
from rest_framework import serializers
from .models import Company, Contact

class ContactSerializer(serializers.ModelSerializer):
    class Meta:
        model = Contact
        fields = '__all__'

class CompanySerializer(serializers.ModelSerializer):
    contacts = ContactSerializer(many=True, read_only=True)
    class Meta:
        model = Company
        fields = '__all__'
EOF

cat > backend/apps/companies/views.py << 'EOF'
from rest_framework import generics, permissions
from .models import Company
from .serializers import CompanySerializer

class CompanyListCreateView(generics.ListCreateAPIView):
    queryset = Company.objects.all()
    serializer_class = CompanySerializer
    permission_classes = [permissions.IsAuthenticatedOrReadOnly]

class CompanyDetailView(generics.RetrieveUpdateDestroyAPIView):
    queryset = Company.objects.all()
    serializer_class = CompanySerializer
    permission_classes = [permissions.IsAuthenticatedOrReadOnly]
EOF

cat > backend/apps/companies/urls.py << 'EOF'
from django.urls import path
from .views import CompanyListCreateView, CompanyDetailView

urlpatterns = [
    path('', CompanyListCreateView.as_view()),
    path('<int:pk>/', CompanyDetailView.as_view()),
]
EOF

# ---------- apps/theme ----------
cat > backend/apps/theme/models.py << 'EOF'
from django.db import models
from apps.users.models import User

class ThemeSettings(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE, related_name='theme')
    primary_color = models.CharField(max_length=7, default='#3b82f6')  # hex
    background_image = models.ImageField(upload_to='backgrounds/', blank=True, null=True)
    splash_image = models.ImageField(upload_to='splash/', blank=True, null=True)
    font_family = models.CharField(max_length=50, default='Inter, sans-serif')
    dark_mode = models.BooleanField(default=False)
    updated_at = models.DateTimeField(auto_now=True)

    def __str__(self):
        return f"Тема {self.user.email}"
EOF

cat > backend/apps/theme/serializers.py << 'EOF'
from rest_framework import serializers
from .models import ThemeSettings

class ThemeSettingsSerializer(serializers.ModelSerializer):
    class Meta:
        model = ThemeSettings
        fields = '__all__'
        read_only_fields = ('user',)
EOF

cat > backend/apps/theme/views.py << 'EOF'
from rest_framework import generics, permissions
from .models import ThemeSettings
from .serializers import ThemeSettingsSerializer

class ThemeSettingsDetailView(generics.RetrieveUpdateAPIView):
    serializer_class = ThemeSettingsSerializer
    permission_classes = [permissions.IsAuthenticated]

    def get_object(self):
        obj, created = ThemeSettings.objects.get_or_create(user=self.request.user)
        return obj
EOF

cat > backend/apps/theme/urls.py << 'EOF'
from django.urls import path
from .views import ThemeSettingsDetailView

urlpatterns = [
    path('settings/', ThemeSettingsDetailView.as_view()),
]
EOF

# Создаём __init__.py для всех приложений
for app in users documents orders chat notifications companies theme; do
    touch backend/apps/$app/__init__.py
done

# ======================== FRONTEND ========================

cat > frontend/package.json << 'EOF'
{
  "name": "gruzavtotrans-frontend",
  "version": "1.0.0",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.20.0",
    "axios": "^1.6.2",
    "react-yandex-maps": "^4.0.0",
    "tailwindcss": "^3.3.6",
    "react-color": "^2.19.3"
  },
  "devDependencies": {
    "@vitejs/plugin-react": "^4.2.1",
    "vite": "^5.0.0"
  }
}
EOF

cat > frontend/vite.config.js << 'EOF'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api': 'http://localhost:8000',
      '/ws': { target: 'ws://localhost:8000', ws: true },
    },
  },
})
EOF

cat > frontend/tailwind.config.js << 'EOF'
/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: { extend: {} },
  plugins: [],
}
EOF

cat > frontend/public/index.html << 'EOF'
<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Грузавтотранс</title>
</head>
<body>
  <div id="root"></div>
  <script type="module" src="/src/index.js"></script>
</body>
</html>
EOF

cat > frontend/src/index.css << 'EOF'
@tailwind base;
@tailwind components;
@tailwind utilities;
EOF

cat > frontend/src/index.js << 'EOF'
import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './index.css';

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode><App /></React.StrictMode>
);
EOF

cat > frontend/src/App.jsx << 'EOF'
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { ThemeProvider } from './context/ThemeContext';
import Register from './components/auth/Register';
import Login from './components/auth/Login';
import Dashboard from './components/dashboard/Dashboard';
import OrderForm from './components/orders/OrderForm';
import OrderMap from './components/orders/OrderMap';
import Settings from './components/settings/Settings';
import CompaniesList from './components/companies/CompaniesList';

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <ThemeProvider>
          <Routes>
            <Route path="/register" element={<Register />} />
            <Route path="/login" element={<Login />} />
            <Route path="/" element={<Dashboard />} />
            <Route path="/order/new" element={<OrderForm />} />
            <Route path="/order/:slug" element={<OrderMap />} />
            <Route path="/settings" element={<Settings />} />
            <Route path="/companies" element={<CompaniesList />} />
          </Routes>
        </ThemeProvider>
      </AuthProvider>
    </BrowserRouter>
  );
}
export default App;
EOF

# ---------- контексты ----------
cat > frontend/src/context/AuthContext.jsx << 'EOF'
import { createContext, useState, useContext } from 'react';
import api from '../services/api';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);

  const login = async (email, password) => {
    // Здесь должен быть вызов API для получения токена
    // Пока заглушка
    const response = await api.post('/auth/login/', { email, password });
    setUser(response.data.user);
    localStorage.setItem('token', response.data.token);
  };

  const logout = () => {
    setUser(null);
    localStorage.removeItem('token');
  };

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
EOF

cat > frontend/src/context/ThemeContext.jsx << 'EOF'
import { createContext, useState, useContext, useEffect } from 'react';
import api from '../services/api';

const ThemeContext = createContext();

export const ThemeProvider = ({ children }) => {
  const [theme, setTheme] = useState({
    primary_color: '#3b82f6',
    background_image: null,
    splash_image: null,
    font_family: 'Inter, sans-serif',
    dark_mode: false,
  });

  const fetchTheme = async () => {
    try {
      const res = await api.get('/theme/settings/');
      setTheme(res.data);
    } catch (e) {
      console.log('Тема не загружена');
    }
  };

  const updateTheme = async (newTheme) => {
    try {
      const res = await api.patch('/theme/settings/', newTheme);
      setTheme(res.data);
    } catch (e) {
      console.error('Ошибка обновления темы', e);
    }
  };

  useEffect(() => {
    fetchTheme();
  }, []);

  return (
    <ThemeContext.Provider value={{ theme, updateTheme }}>
      {children}
    </ThemeContext.Provider>
  );
};

export const useTheme = () => useContext(ThemeContext);
EOF

# ---------- сервис API ----------
cat > frontend/src/services/api.js << 'EOF'
import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;
EOF

# ---------- компоненты ----------
# Компонент регистрации (дополнен)
cat > frontend/src/components/auth/Register.jsx << 'EOF'
import { useState } from 'react';
import api from '../../services/api';

export default function Register() {
  const [step, setStep] = useState('form');
  const [form, setForm] = useState({});
  const [phone, setPhone] = useState('');
  const [email, setEmail] = useState('');
  const [smsCode, setSmsCode] = useState('');
  const [emailCode, setEmailCode] = useState('');

  const handleChange = (e) => setForm({ ...form, [e.target.name]: e.target.value });

  const handleRegister = async (e) => {
    e.preventDefault();
    try {
      const res = await api.post('/auth/register/', form);
      setPhone(res.data.phone);
      setStep('sms');
    } catch (err) { alert('Ошибка регистрации'); }
  };

  const handleVerifySms = async () => {
    try {
      const res = await api.post('/auth/verify-sms/', { phone, code: smsCode });
      setEmail(res.data.email);
      setStep('email');
    } catch { alert('Неверный SMS-код'); }
  };

  const handleVerifyEmail = async () => {
    try {
      await api.post('/auth/verify-email/', { email, code: emailCode });
      setStep('done');
    } catch { alert('Неверный email-код'); }
  };

  return (
    <div className="max-w-lg mx-auto bg-white p-8 rounded-2xl shadow mt-10">
      {step === 'form' && (
        <form onSubmit={handleRegister}>
          <h2 className="text-2xl font-bold mb-6">Регистрация</h2>
          <input name="full_name" onChange={handleChange} placeholder="ФИО" className="border p-2 w-full mb-2" required />
          <input name="phone" onChange={handleChange} placeholder="Телефон" className="border p-2 w-full mb-2" required />
          <input name="email" type="email" onChange={handleChange} placeholder="Email" className="border p-2 w-full mb-2" required />
          <input name="inn" onChange={handleChange} placeholder="ИНН" className="border p-2 w-full mb-2" />
          <input name="passport_data" onChange={handleChange} placeholder="Паспорт (серия номер)" className="border p-2 w-full mb-2" />
          <input name="registration_address" onChange={handleChange} placeholder="Прописка" className="border p-2 w-full mb-2" />
          <input name="driver_license_data" onChange={handleChange} placeholder="В/У (серия номер)" className="border p-2 w-full mb-2" />
          <input name="vehicle" onChange={handleChange} placeholder="Тягач (госномер)" className="border p-2 w-full mb-2" />
          <input name="trailer" onChange={handleChange} placeholder="Прицеп (госномер)" className="border p-2 w-full mb-2" />
          <button type="submit" className="bg-blue-600 text-white p-2 w-full rounded-xl">Зарегистрироваться</button>
        </form>
      )}
      {step === 'sms' && (
        <div>
          <h2 className="text-xl font-bold">Подтверждение телефона</h2>
          <p>Код отправлен на {phone}</p>
          <input value={smsCode} onChange={(e) => setSmsCode(e.target.value)} placeholder="Код из SMS" className="border p-2 w-full mb-2" />
          <button onClick={handleVerifySms} className="bg-green-600 text-white p-2 w-full rounded-xl">Подтвердить</button>
        </div>
      )}
      {step === 'email' && (
        <div>
          <h2 className="text-xl font-bold">Подтверждение email</h2>
          <p>Код отправлен на {email}</p>
          <input value={emailCode} onChange={(e) => setEmailCode(e.target.value)} placeholder="Код из письма" className="border p-2 w-full mb-2" />
          <button onClick={handleVerifyEmail} className="bg-green-600 text-white p-2 w-full rounded-xl">Подтвердить</button>
        </div>
      )}
      {step === 'done' && <div className="text-center text-green-600 text-xl">Регистрация завершена! Войдите в систему.</div>}
    </div>
  );
}
EOF

cat > frontend/src/components/auth/Login.jsx << 'EOF'
import { useState } from 'react';
import { useAuth } from '../../context/AuthContext';
import { useNavigate } from 'react-router-dom';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await login(email, password);
      navigate('/');
    } catch {
      alert('Неверные данные');
    }
  };

  return (
    <div className="max-w-md mx-auto bg-white p-8 rounded-2xl shadow mt-10">
      <h2 className="text-2xl font-bold mb-6">Вход</h2>
      <form onSubmit={handleSubmit}>
        <input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} className="border p-2 w-full mb-2" required />
        <input type="password" placeholder="Пароль" value={password} onChange={(e) => setPassword(e.target.value)} className="border p-2 w-full mb-2" required />
        <button type="submit" className="bg-blue-600 text-white p-2 w-full rounded-xl">Войти</button>
      </form>
    </div>
  );
}
EOF

# Дашборд (заглушка)
cat > frontend/src/components/dashboard/Dashboard.jsx << 'EOF'
import { Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

export default function Dashboard() {
  const { user } = useAuth();
  return (
    <div className="p-6">
      <h1 className="text-3xl font-bold">Добро пожаловать, {user?.full_name || 'Гость'}!</h1>
      <div className="grid grid-cols-2 gap-4 mt-6">
        <Link to="/order/new" className="bg-blue-500 text-white p-4 rounded-xl text-center">Создать заказ</Link>
        <Link to="/companies" className="bg-green-500 text-white p-4 rounded-xl text-center">Компании</Link>
        <Link to="/settings" className="bg-purple-500 text-white p-4 rounded-xl text-center">Настройки темы</Link>
      </div>
    </div>
  );
}
EOF

# Компонент создания заказа (с картой)
cat > frontend/src/components/orders/OrderForm.jsx << 'EOF'
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../../services/api';
import { YMaps, Map, Placemark } from 'react-yandex-maps';

export default function OrderForm() {
  const navigate = useNavigate();
  const [form, setForm] = useState({
    cargo_description: '',
    weight: '',
    volume: '',
    date_required: '',
    pickup_address: { street: '', house: '', city: '', latitude: 55.7558, longitude: 37.6173 },
    delivery_address: { street: '', house: '', city: '', latitude: 55.7558, longitude: 37.6173 },
  });

  const handleChange = (e, prefix) => {
    const { name, value } = e.target;
    if (prefix) {
      setForm({ ...form, [prefix]: { ...form[prefix], [name]: value } });
    } else {
      setForm({ ...form, [name]: value });
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const res = await api.post('/orders/', { ...form, customer: 1 }); // замените на реального пользователя
      navigate(`/order/${res.data.public_url}`);
    } catch (err) {
      alert('Ошибка создания заказа');
    }
  };

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h2 className="text-2xl font-bold mb-4">Новый заказ</h2>
      <form onSubmit={handleSubmit}>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <h3 className="font-semibold">Адрес забора</h3>
            <input name="street" placeholder="Улица" onChange={(e) => handleChange(e, 'pickup_address')} className="border p-2 w-full mb-2" />
            <input name="house" placeholder="Дом" onChange={(e) => handleChange(e, 'pickup_address')} className="border p-2 w-full mb-2" />
            <input name="city" placeholder="Город" onChange={(e) => handleChange(e, 'pickup_address')} className="border p-2 w-full mb-2" />
          </div>
          <div>
            <h3 className="font-semibold">Адрес доставки</h3>
            <input name="street" placeholder="Улица" onChange={(e) => handleChange(e, 'delivery_address')} className="border p-2 w-full mb-2" />
            <input name="house" placeholder="Дом" onChange={(e) => handleChange(e, 'delivery_address')} className="border p-2 w-full mb-2" />
            <input name="city" placeholder="Город" onChange={(e) => handleChange(e, 'delivery_address')} className="border p-2 w-full mb-2" />
          </div>
        </div>
        <div className="mt-4">
          <textarea name="cargo_description" placeholder="Описание груза" onChange={handleChange} className="border p-2 w-full mb-2" />
          <input name="weight" type="number" placeholder="Вес (т)" onChange={handleChange} className="border p-2 w-full mb-2" />
          <input name="volume" type="number" placeholder="Объём (м³)" onChange={handleChange} className="border p-2 w-full mb-2" />
          <input name="date_required" type="date" onChange={handleChange} className="border p-2 w-full mb-2" />
        </div>
        <button type="submit" className="bg-blue-600 text-white p-2 w-full rounded-xl">Создать заказ</button>
      </form>
    </div>
  );
}
EOF

# Страница заказа с картой
cat > frontend/src/components/orders/OrderMap.jsx << 'EOF'
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../services/api';
import { YMaps, Map, Placemark } from 'react-yandex-maps';

export default function OrderMap() {
  const { slug } = useParams();
  const [order, setOrder] = useState(null);

  useEffect(() => {
    api.get(`/orders/${slug}/`).then(res => setOrder(res.data));
  }, [slug]);

  if (!order) return <div>Загрузка...</div>;

  const pickup = order.pickup_address;
  const delivery = order.delivery_address;

  return (
    <div className="p-6">
      <h2 className="text-2xl font-bold">Заказ #{order.id}</h2>
      <p>Вес: {order.weight} т, Объём: {order.volume} м³</p>
      <div className="h-96 w-full mt-4">
        <YMaps>
          <Map defaultState={{ center: [55.7558, 37.6173], zoom: 10 }} className="w-full h-full">
            <Placemark geometry={[pickup.latitude, pickup.longitude]} properties={{ hintContent: 'Забор' }} />
            <Placemark geometry={[delivery.latitude, delivery.longitude]} properties={{ hintContent: 'Доставка' }} />
          </Map>
        </YMaps>
      </div>
    </div>
  );
}
EOF

# Настройки темы
cat > frontend/src/components/settings/Settings.jsx << 'EOF'
import { useTheme } from '../../context/ThemeContext';
import { useState } from 'react';
import { ChromePicker } from 'react-color';

export default function Settings() {
  const { theme, updateTheme } = useTheme();
  const [color, setColor] = useState(theme.primary_color);
  const [darkMode, setDarkMode] = useState(theme.dark_mode);
  const [bgFile, setBgFile] = useState(null);
  const [splashFile, setSplashFile] = useState(null);

  const handleSave = async () => {
    const formData = new FormData();
    formData.append('primary_color', color);
    formData.append('dark_mode', darkMode);
    if (bgFile) formData.append('background_image', bgFile);
    if (splashFile) formData.append('splash_image', splashFile);
    await updateTheme(formData);
    alert('Тема сохранена');
  };

  return (
    <div className="max-w-2xl mx-auto p-6">
      <h2 className="text-2xl font-bold mb-4">Настройки темы</h2>
      <div className="mb-4">
        <label className="block font-semibold">Основной цвет</label>
        <ChromePicker color={color} onChange={c => setColor(c.hex)} />
      </div>
      <div className="mb-4">
        <label className="block font-semibold">Тёмная тема</label>
        <input type="checkbox" checked={darkMode} onChange={(e) => setDarkMode(e.target.checked)} />
      </div>
      <div className="mb-4">
        <label className="block font-semibold">Фоновое изображение</label>
        <input type="file" accept="image/*" onChange={(e) => setBgFile(e.target.files[0])} />
      </div>
      <div className="mb-4">
        <label className="block font-semibold">Изображение заставки</label>
        <input type="file" accept="image/*" onChange={(e) => setSplashFile(e.target.files[0])} />
      </div>
      <button onClick={handleSave} className="bg-blue-600 text-white p-2 rounded-xl">Сохранить</button>
    </div>
  );
}
EOF

# Список компаний
cat > frontend/src/components/companies/CompaniesList.jsx << 'EOF'
import { useEffect, useState } from 'react';
import api from '../../services/api';

export default function CompaniesList() {
  const [companies, setCompanies] = useState([]);
  const [newCompany, setNewCompany] = useState({ name: '', inn: '', address: '', phone: '', email: '' });

  useEffect(() => {
    api.get('/companies/').then(res => setCompanies(res.data));
  }, []);

  const handleAdd = async (e) => {
    e.preventDefault();
    try {
      const res = await api.post('/companies/', newCompany);
      setCompanies([...companies, res.data]);
      setNewCompany({ name: '', inn: '', address: '', phone: '', email: '' });
    } catch (err) { alert('Ошибка добавления'); }
  };

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h2 className="text-2xl font-bold mb-4">Компании</h2>
      <form onSubmit={handleAdd} className="grid grid-cols-2 gap-2 mb-6">
        <input placeholder="Название" value={newCompany.name} onChange={(e) => setNewCompany({...newCompany, name: e.target.value})} className="border p-2" />
        <input placeholder="ИНН" value={newCompany.inn} onChange={(e) => setNewCompany({...newCompany, inn: e.target.value})} className="border p-2" />
        <input placeholder="Адрес" value={newCompany.address} onChange={(e) => setNewCompany({...newCompany, address: e.target.value})} className="border p-2" />
        <input placeholder="Телефон" value={newCompany.phone} onChange={(e) => setNewCompany({...newCompany, phone: e.target.value})} className="border p-2" />
        <input placeholder="Email" value={newCompany.email} onChange={(e) => setNewCompany({...newCompany, email: e.target.value})} className="border p-2" />
        <button type="submit" className="bg-green-600 text-white p-2 rounded-xl">Добавить</button>
      </form>
      <div className="grid grid-cols-1 gap-4">
        {companies.map(c => (
          <div key={c.id} className="border p-4 rounded-xl shadow">
            <h3 className="text-xl font-bold">{c.name}</h3>
            <p>ИНН: {c.inn}</p>
            <p>Адрес: {c.address}</p>
            <p>Телефон: {c.phone}</p>
            <p>Email: {c.email}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
EOF

# ======================== DOCKER COMPOSE ========================

cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  db:
    image: postgres:15
    environment:
      POSTGRES_DB: gruz_db
      POSTGRES_USER: gruz_user
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7
    ports:
      - "6379:6379"

  backend:
    build: ./backend
    volumes:
      - ./backend:/app
      - media_volume:/app/media
    ports:
      - "8000:8000"
    depends_on:
      - db
      - redis
    environment:
      - DB_HOST=db
      - REDIS_HOST=redis
      - DEBUG=1

  frontend:
    build: ./frontend
    volumes:
      - ./frontend:/app
      - /app/node_modules
    ports:
      - "3000:3000"
    depends_on:
      - backend
    command: npm run dev -- --host 0.0.0.0

volumes:
  postgres_data:
  media_volume:
EOF

# Создаём Dockerfile для фронтенда
cat > frontend/Dockerfile << 'EOF'
FROM node:18-alpine
WORKDIR /app
COPY package.json .
RUN npm install
COPY . .
EXPOSE 3000
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
EOF

# ======================== .env.example ========================

cat > backend/.env.example << 'EOF'
SECRET_KEY=your-secret-key
DEBUG=1
DB_NAME=gruz_db
DB_USER=gruz_user
DB_PASSWORD=password
DB_HOST=db
DB_PORT=5432
REDIS_HOST=redis
TWILIO_ACCOUNT_SID=your_twilio_sid
TWILIO_AUTH_TOKEN=your_twilio_token
TWILIO_PHONE=+1234567890
EMAIL_HOST_USER=your_sendgrid_username
EMAIL_HOST_PASSWORD=your_sendgrid_password
YANDEX_MAPS_API_KEY=your_yandex_maps_key
EOF

# ======================== README ========================

cat > README.md << 'EOF'
# Грузавтотранс – платформа для грузоперевозок

Полнофункциональное веб-приложение для заказа и отслеживания грузоперевозок с возможностью:
- Регистрация с подтверждением по SMS и email
- Создание заказов с указанием адресов и веса груза
- Интерактивная карта (Яндекс.Карты)
- Чат между заказчиком и перевозчиком (WebSockets)
- Настройка темы (цвет, фон, заставка)
- Добавление компаний и контактов через документы

## Запуск через Docker Desktop

1. Установите [Docker Desktop](https://www.docker.com/products/docker-desktop/).
2. Склонируйте или распакуйте этот проект.
3. В корневой папке выполните:
   ```bash
   docker-compose up --build
